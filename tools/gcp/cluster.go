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
	"os"
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
	clusterName string,
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
		ClusterName: clusterName,
		Config:      clusterConfig,
	}

	request := &dataprocpb.CreateClusterRequest{
		ProjectId: projectName,
		Cluster:   cluster,
		Region:    gcpRegion,
	}

	return request
}

// Register and request parsing of cli parameters with Donna.
func initClusterCmd() {
	donna.ExpectStrOption("b", "bucket", "Bucket name, for cluster setup.", "")
	donna.ExpectStrOption("r", "cred", "Path to credential file.", "")
	donna.ExpectStrOption("n", "name", "Name of the cluster.", "")
	donna.ExpectIntOption("c", "cores", "Number of cores.", 0)
	donna.ExpectIntOption("w", "numWorkers", "Number of workers.", 0)
	donna.ExpectFlag("m", "highmem", "Indicates whether high memory instances should be used.")
	err := donna.Parse()
	if err != nil {
		fmt.Println(err.Error())
		donna.DisplayCommandHelp()
	}
}

// Parses and returns project ID, credential and staging bucket.
func getGcpVars() (string, string, string) {
	// Project ID is inferred from the credential path. The credential path is
	// possibly inferred by the bucket.  Thus the parsing order must  be bucket,
	// credential path and project ID, respectively.

	bucketName, passed := donna.GetStrOption("bucket")
	if !passed {
		fmt.Println("Did not receive bucket name.")
		donna.DisplayCommandHelp()
	}

	credentialPath := getCredentials(bucketName)

	projectName, err := getProjectId(credentialPath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return projectName, credentialPath, bucketName
}

// Parses cluster name, number of workers, number of cores and high memory flag.
func getClusterVars() (string, int, int, bool) {
	clusterName, passed := donna.GetStrOption("name")
	if !passed {
		clusterName = "ro-cluster"
	}

	workerCount, passed := donna.GetIntOption("workers")
	if !passed {
		fmt.Println("Did not receive number of workers.")
		donna.DisplayCommandHelp()
	}

	coreCount, passed := donna.GetIntOption("cores")
	if !passed {
		fmt.Println("Did not receive number of cores.")
		donna.DisplayCommandHelp()
	}

	highMemory := donna.HasFlag("highmem")

	return clusterName, workerCount, coreCount, highMemory
}

func clusterCmd() {
	// Register and request parsing of cli parameters.
	initClusterCmd()

	// Load and validate cli parameters.
	projectId, credentialPath, bucketName := getGcpVars()
	clusterName, workerCount, coreCount, highMemory := getClusterVars()

	// Create a cluster controller client.
	ctx := context.Background()
	client, err := dataproc.NewClusterControllerClient(
		ctx,
		option.WithCredentialsFile(credentialPath),
		option.WithEndpoint(gcpEndpoint),
	)
	if err != nil {
		logger.Fatal("Failed to create a cluster controller client: %s", err)
	}
	defer client.Close()

	// Compute appropriate values and generate a request for cluster creation.
	request := getClusterCreationRequest(
		clusterName,
		projectId,
		bucketName,
		workerCount,
		coreCount,
		highMemory,
	)

	// Perform actual request to create the cluster.
	_, err = client.CreateCluster(ctx, request)
	if err != nil {
		logger.Fatal("Could not create cluster: %s", err)
	}

	fmt.Println("Cluster creation request successful.")
}
