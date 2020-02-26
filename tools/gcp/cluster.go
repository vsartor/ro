// ro - Copyright (c) Victhor Sartório, 2020-. All rights reserved.
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package gcp

import (
	dataproc "cloud.google.com/go/dataproc/apiv1beta2"
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/vsartor/ro/donna"
	"google.golang.org/api/option"
	dataprocpb "google.golang.org/genproto/googleapis/cloud/dataproc/v1beta2"
	"math"
	"strconv"
)

const (
	workerOsMemory = 4096
	overheadRatio  = 0.1
)

const (
	gcpRegion   = "us-central1"
	gcpZone     = "us-central1-b"
	gcpEndpoint = "us-central1-dataproc.googleapis.com:443"
)

// Computes an appropriate value for yarn.nodemanager.resource.memory-mb
// setting based on on total node machine capacity.
func yarnMemory(numCores, memoryPerCore int) int {
	totalMemory := numCores * memoryPerCore
	return totalMemory - workerOsMemory
}

// Computes an appropriate value for spark.executor.memory setting based
// on number of cores per executor and amount of memory per core in the
// machine.
func computeExecutorMemory(cores, memoryPerCore int) int {
	exactMemoryValue := (1.0 - overheadRatio) * float64(yarnMemory(cores, memoryPerCore))
	return int(math.Floor(exactMemoryValue))
}

// Builds a CreateClusterRequest to be send by the client
func getClusterCreationRequest(
	projectName string,
	bucketName string,
	numWorkers int,
	numCores int,
	highMemory bool,
) *dataprocpb.CreateClusterRequest {
	// Compute the tier of the machines we'll work with
	machineUriName := "n1-standard"
	if highMemory {
		machineUriName = "n1-highmem"
	}

	// Compute master machine configuration
	masterType := "n1-standard-2"
	masterMemory := 4096
	masterMaxResult := 3072
	if numWorkers >= 6 || highMemory {
		masterType = "n1-highmem-4"
		masterMemory = 18432
		masterMaxResult = 16384
	}

	// Compute worker machine configuration
	workerType := fmt.Sprintf("%s-%d", machineUriName, numCores)
	memoryPerCore := 3840
	if highMemory {
		memoryPerCore = 6656
	}
	memoryPerWorker := yarnMemory(numCores, memoryPerCore)
	executorMemory := computeExecutorMemory(numCores, memoryPerCore)

	// Below here we just build the request

	diskConfig := &dataprocpb.DiskConfig{
		BootDiskType:   "pd-standard",
		BootDiskSizeGb: 1000,
	}

	masterConfig := &dataprocpb.InstanceGroupConfig{
		NumInstances:   1,
		MachineTypeUri: masterType,
		DiskConfig:     diskConfig,
	}

	workerConfig := &dataprocpb.InstanceGroupConfig{
		NumInstances:   int32(numWorkers),
		MachineTypeUri: workerType,
		DiskConfig:     diskConfig,
	}

	endpointConfig := &dataprocpb.EndpointConfig{
		EnableHttpPortAccess: true,
	}

	gceClusterConfig := &dataprocpb.GceClusterConfig{
		ZoneUri: gcpZone,
	}

	softwareConfig := &dataprocpb.SoftwareConfig{
		ImageVersion: "1.4",
		OptionalComponents: []dataprocpb.Component{
			dataprocpb.Component_ANACONDA,
			dataprocpb.Component_JUPYTER,
		},
		Properties: map[string]string{
			"core:fs.gs.implicit.dir.repair.enable":     "false",
			"core:fs.gs.requester.pays.mode":            "AUTO",
			"core:fs.gs.requester.pays.project.id":      projectName,
			"spark:spark.driver.maxResultSize":          fmt.Sprintf("%dm", masterMaxResult),
			"spark:spark.driver.memory":                 fmt.Sprintf("%dm", masterMemory),
			"spark:spark.executor.instances":            fmt.Sprintf("%d", numWorkers),
			"spark:spark.executor.cores":                fmt.Sprintf("%d", numCores),
			"spark:spark.executor.memory":               fmt.Sprintf("%dm", executorMemory),
			"spark:spark.jars.packages":                 "org.apache.spark:spark-avro_2.11:2.4.4",
			"spark:spark.dynamicAllocation.enabled":     "false",
			"yarn:yarn.nodemanager.resource.memory-mb":  fmt.Sprintf("%d", memoryPerWorker),
			"yarn:yarn.scheduler.maximum-allocation-mb": fmt.Sprintf("%d", memoryPerWorker),
		},
	}

	lifecycleConfig := &dataprocpb.LifecycleConfig{
		IdleDeleteTtl: &duration.Duration{Seconds: 3600},
	}

	clusterConfig := &dataprocpb.ClusterConfig{
		ConfigBucket:     bucketName,
		GceClusterConfig: gceClusterConfig,
		MasterConfig:     masterConfig,
		WorkerConfig:     workerConfig,
		SoftwareConfig:   softwareConfig,
		LifecycleConfig:  lifecycleConfig,
		EndpointConfig:   endpointConfig,
	}

	cluster := &dataprocpb.Cluster{
		ProjectId:   projectName,
		ClusterName: "ro-cluster",
		Config:      clusterConfig,
	}

	request := &dataprocpb.CreateClusterRequest{
		ProjectId: projectName,
		Cluster:   cluster,
		Region:    gcpRegion,
	}

	return request
}

func clusterCmd() {
	// Handle command line arguments

	donna.ExpectOption("project")
	donna.ExpectOption("bucket")
	donna.ExpectOption("cred")
	donna.ExpectOption("c")
	donna.ExpectOption("w")
	donna.ExpectFlag("highmem")
	donna.ValidateLocal()

	projectName, ok := donna.GetOption("project")
	if !ok {
		logger.Fatal("Did not receive a project name.")
	}
	bucketName, ok := donna.GetOption("bucket")
	if !ok {
		logger.Fatal("Did not receive a bucket name.")
	}
	credential, ok := donna.GetOption("cred")
	if !ok {
		logger.Fatal("Did not receive path to credential file.")
	}
	cores, ok := donna.GetOption("c")
	if !ok {
		logger.Fatal("Did not receive number of cores.")
	}
	workers, ok := donna.GetOption("w")
	if !ok {
		logger.Fatal("Did not receive number of workers.")
	}
	highMemory := donna.HasFlag("highmem")

	// Convert integer options
	numWorkers, err := strconv.Atoi(workers)
	if err != nil {
		logger.Fatal("'%s' is not a valid number of workers.", workers)
	}
	numCores, err := strconv.Atoi(cores)
	if err != nil {
		logger.Fatal("'%s' is not a valid number of cores.", cores)
	}

	// Create the cluster

	ctx := context.Background()
	client, err := dataproc.NewClusterControllerClient(
		ctx,
		option.WithCredentialsFile(credential),
		option.WithEndpoint(gcpEndpoint),
	)
	if err != nil {
		logger.Fatal("Failed to create clusterControllerClient: %s", err)
	}
	defer client.Close()

	request := getClusterCreationRequest(projectName, bucketName, numWorkers, numCores, highMemory)

	_, err = client.CreateCluster(ctx, request)
	if err != nil {
		logger.Fatal("Could not create cluster: %s", err)
	}

	logger.Info("Cluster creation request successfully sent.")
}
