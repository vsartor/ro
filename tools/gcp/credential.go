// ro - Copyright (c) Victhor Sartório, 2020-. All rights reserved.
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package gcp

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vsartor/ro/bran"
	"github.com/vsartor/ro/donna"
	"io/ioutil"
)

// Returns the path to the credential file.
// If the user passed the `r/cred` option, the value is fetched from here
// and this credential will be associated by Bran to the bucket path.
// If the credential was not passed but Bran remembers an associated credential
// to this bucket, the value remembered by Bran will be used.
// If neither the user passed nor Bran remembers a valid credential, an error
// message is printed and the program exits by displaying the usage information
// for whatever command invoked this routine.
func getCredentials(bucket string) string {
	key := fmt.Sprintf("cred:bucket:%s", bucket)

	// If a credential was passed, remember it and use it
	credential, passed := donna.GetStrOption("cred")
	if passed {
		err := bran.Set(key, credential)
		logger.Info("Associating credential %q to bucket %q", credential, bucket)
		if err != nil {
			panic(err.Error())
		}
		return credential
	}

	// If credential wasn't passed try to fetch it from memory
	credential, err := bran.Get(key)
	if err != nil {
		panic(err.Error())
	}

	if credential == "" {
		fmt.Println("Credential wasn't passed and couldn't be remembered.")
		donna.DisplayCommandHelp()
	}

	return credential
}

// Returns the project ID associated with a credential file.
func getProjectId(credentialPath string) (string, error) {
	credentialContents, err := ioutil.ReadFile(credentialPath)
	if err != nil {
		return "", err
	}

	credential := make(map[string]string)
	err = json.Unmarshal(credentialContents, &credential)
	if err != nil {
		return "", err
	}

	projectId, ok := credential["project_id"]
	if !ok {
		return "", errors.New("JSON credential missing `project_id` field")
	}

	return projectId, nil
}
