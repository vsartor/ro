package donna

import (
	"errors"
	"fmt"
)

// Type indicating the type of a parameter.
type ParamKind uint

const (
	ParamStr ParamKind = iota
	ParamFlag
	ParamInt
)

// Holds parameter information built by the parser.
type ParamInfo struct {
	name      string    // Parameter name.
	kind      ParamKind // Parameter type.
	valueStr  string    // Parameter value, passed or default.
	valueInt  int
	valueBool bool
	wasPassed bool      // Indicates whether parameter was passed.
}

func NewParamInfo(expectedParam ParamExpectInfo) ParamInfo {
	return ParamInfo{
		name:      expectedParam.name,
		kind:      expectedParam.kind,
		valueStr:  expectedParam.defaultStr,
		valueInt:  expectedParam.defaultInt,
		valueBool: false,
		wasPassed: false,
	}
}

func (paramInfo *ParamInfo) ToggleFlag() {
	paramInfo.valueBool = true
}

func (paramInfo *ParamInfo) SetStrValue(value string) {
	paramInfo.valueStr = value
	paramInfo.wasPassed = true
}

func (paramInfo *ParamInfo) SetIntValue(value int) {
	paramInfo.valueInt = value
	paramInfo.wasPassed = true
}

var (
	globalParams map[string]ParamInfo
	localParams  map[string]ParamInfo
)

func init() {
	globalParams = make(map[string]ParamInfo)
	localParams = make(map[string]ParamInfo)
}

// Implements common functionality for checking the presence of
// a flag for HasGlobalFlag and HasLocalFlag.
func hasFlag(name string, params map[string]ParamInfo) (bool, error) {
	info, ok := params[name]
	if !ok {
		errorMsg := fmt.Sprintf("Unknown parameter name '%s'.", name)
		return false, errors.New(errorMsg)
	}

	if info.kind != ParamFlag {
		errorMsg := fmt.Sprintf("Parameter '%s' is not a flag.", name)
		return false, errors.New(errorMsg)
	}

	return info.valueBool, nil
}

// Returns whether the flag was passed as a global parameter.
// If parameter name is incorrect, this is assumed to be a developer
// error and thus a panic is thrown to aid in a quick fix.
func HasGlobalFlag(name string) bool {
	returnVal, err := hasFlag(name, globalParams)
	if err != nil {
		panic(err.Error())
	}
	return returnVal
}

// Returns whether the flag was passed as a local parameter.
// If parameter name is incorrect, this is assumed to be a developer
// error and thus a panic is thrown to aid in a quick fix.
func HasFlag(name string) bool {
	returnVal, err := hasFlag(name, localParams)
	if err != nil {
		panic(err.Error())
	}
	return returnVal
}

// Implements common functionality for obtaining a string option value.
func getStrOption(name string, params map[string]ParamInfo) (string, bool, error) {
	info, ok := params[name]
	if !ok {
		errorMsg := fmt.Sprintf("Unknown parameter name '%s'.", name)
		return "", false, errors.New(errorMsg)
	}

	if info.kind != ParamStr {
		errorMsg := fmt.Sprintf("Parameter '%s' is not a string option.", name)
		return "", false, errors.New(errorMsg)
	}

	return info.valueStr, info.wasPassed, nil
}

// Returns the global string option's value and whether value was passed.
// If parameter name is incorrect, this is assumed to be a developer
// error and thus a panic is thrown to aid in a quick fix.
func GetGlobalStrOption(name string) (string, bool) {
	returnVal, ok, err := getStrOption(name, globalParams)
	if err != nil {
		panic(err.Error())
	}
	return returnVal, ok
}

// Returns the local string option's value and whether value was passed.
// If parameter name is incorrect, this is assumed to be a developer
// error and thus a panic is thrown to aid in a quick fix.
func GetStrOption(name string) (string, bool) {
	returnVal, ok, err := getStrOption(name, localParams)
	if err != nil {
		panic(err.Error())
	}
	return returnVal, ok
}

// Implements common functionality for obtaining an integer option value.
func getIntOption(name string, params map[string]ParamInfo) (int, bool, error) {
	info, ok := params[name]
	if !ok {
		errorMsg := fmt.Sprintf("Unknown parameter name '%s'.", name)
		return 0, false, errors.New(errorMsg)
	}

	if info.kind != ParamInt {
		errorMsg := fmt.Sprintf("Parameter '%s' is not an integer option.", name)
		return 0, false, errors.New(errorMsg)
	}

	return info.valueInt, info.wasPassed, nil
}

// Returns the global integer option's value and whether value was passed.
// If parameter name is incorrect, this is assumed to be a developer
// error and thus a panic is thrown to aid in a quick fix.
func GetGlobalIntOption(name string) (int, bool) {
	value, ok, err := getIntOption(name, globalParams)
	if err != nil {
		panic(err.Error())
	}
	return value, ok
}

// Returns the local string option's value and whether value was passed.
// If parameter name is incorrect, this is assumed to be a developer
// error and thus a panic is thrown to aid in a quick fix.
func GetIntOption(name string) (int, bool) {
	value, ok, err := getIntOption(name, localParams)
	if err != nil {
		panic(err.Error())
	}
	return value, ok
}
