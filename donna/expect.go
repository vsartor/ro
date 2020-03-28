// ro - Copyright (c) Victhor Sartório, 2020-. All rights reserved.
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package donna

// Expected information about parameters held by the validator.
type ParamExpectInfo struct {
	alias      string    // Short alias for the command.
	name       string    // Name of the parameter.
	kind       ParamKind // Type of the parameter.
	defaultStr string    // Default value for parameter, if it's a string.
	defaultInt int       // Default value for parameter, if it's an integer.
	desc       string    // Parameter description.
}

// Expected information about an argument.
type expectedArgInfo struct {
	name string // Alias for referring to this argument.
	desc string // Argument description.
}

var (
	expectedGlobalParams []ParamExpectInfo
	expectedLocalParams  []ParamExpectInfo
	expectedArgs         []expectedArgInfo
)

// Returns the parameter type. Also returns a flag indicating if
// the request was valid.
func expectedInfo(passedName string, global bool) (ParamExpectInfo, bool) {
	var expectedParams []ParamExpectInfo
	if global {
		expectedParams = expectedGlobalParams
	} else {
		expectedParams = expectedLocalParams
	}

	for _, info := range expectedParams {
		if passedName == info.name || passedName == info.alias {
			return info, true
		}
	}

	return ParamExpectInfo{}, false
}

func expectFlag(alias, name, desc string, where *[]ParamExpectInfo) {
	*where = append(
		*where,
		ParamExpectInfo{
			alias:      alias,
			name:       name,
			kind:       ParamFlag,
			defaultStr: "",
			defaultInt: 0,
			desc:       desc,
		},
	)
}

// Registers a global flag.
func ExpectGlobalFlag(alias, name, desc string) {
	expectFlag(alias, name, desc, &expectedGlobalParams)
}

// Registers a command flag.
func ExpectFlag(alias, name, desc string) {
	expectFlag(alias, name, desc, &expectedLocalParams)
}

func expectStrOption(alias, name, desc, defaultValue string, where *[]ParamExpectInfo) {
	*where = append(
		*where,
		ParamExpectInfo{
			alias:      alias,
			name:       name,
			kind:       ParamStr,
			defaultStr: defaultValue,
			defaultInt: 0,
			desc:       desc,
		},
	)
}

// Register a global string option.
func ExpectGlobalStrOption(alias, name, desc, defaultValue string) {
	expectStrOption(alias, name, desc, defaultValue, &expectedGlobalParams)
}

// Registers a command string option.
func ExpectStrOption(alias, name, desc, defaultValue string) {
	expectStrOption(alias, name, desc, defaultValue, &expectedLocalParams)
}

func expectIntOption(alias, name, desc string, defaultValue int, where *[]ParamExpectInfo) {
	*where = append(
		*where,
		ParamExpectInfo{
			alias:      alias,
			name:       name,
			kind:       ParamInt,
			defaultStr: "",
			defaultInt: defaultValue,
			desc:       desc,
		},
	)
}

// Register a global string option.
func ExpectGlobalIntOption(alias, name, desc string, defaultValue int) {
	expectIntOption(alias, name, desc, defaultValue, &expectedGlobalParams)
}

// Registers a command string option.
func ExpectIntOption(alias, name, desc string, defaultValue int) {
	expectIntOption(alias, name, desc, defaultValue, &expectedLocalParams)
}

// Registers an expected argument.
// Only useful for building the help command.
func ExpectArg(name, desc string) {
	expectedArgs = append(
		expectedArgs,
		expectedArgInfo{
			name: name,
			desc: desc,
		},
	)
}
