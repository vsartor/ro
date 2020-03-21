// ro - Copyright (c) Victhor Sartório, 2020-. All rights reserved.
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package gcp

import (
	"fmt"
	"github.com/vsartor/ro/bran"
	"github.com/vsartor/ro/donna"
	"os"
)

// Gets credential.
func getCredential(bucket string) string {
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
		os.Exit(0)
	}

	return credential
}
