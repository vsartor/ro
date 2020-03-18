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
}

func NewParamInfo(expectedParam ParamExpectInfo) ParamInfo {
	return ParamInfo{
		name:      expectedParam.name,
		kind:      expectedParam.kind,
		valueStr:  expectedParam.defaultStr,
		valueInt:  expectedParam.defaultInt,
		valueBool: false,
	}
}

func (paramInfo *ParamInfo) ToggleFlag() {
	paramInfo.valueBool = true
}

func (paramInfo *ParamInfo) SetStrValue(value string) {
	paramInfo.valueStr = value
}

func (paramInfo *ParamInfo) SetIntValue(value int) {
	paramInfo.valueInt = value
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
// Returns an error in case the parameter name was incorrect.
func HasGlobalFlag(name string) (bool, error) {
	return hasFlag(name, globalParams)
}

// Returns whether the flag was passed as a local parameter.
// Returns an error in case the parameter name was incorrect.
func HasFlag(name string) (bool, error) {
	return hasFlag(name, localParams)
}

// Implements common functionality for obtaining a string option value.
func getStrOption(name string, params map[string]ParamInfo) (string, error) {
	info, ok := params[name]
	if !ok {
		errorMsg := fmt.Sprintf("Unknown parameter name '%s'.", name)
		return "", errors.New(errorMsg)
	}

	if info.kind != ParamStr {
		errorMsg := fmt.Sprintf("Parameter '%s' is not a string option.", name)
		return "", errors.New(errorMsg)
	}

	return info.valueStr, nil
}

// Returns the global string option's value.
// Returns an error in case the parameter name was incorrect.
func GetGlobalStrOption(name string) (string, error) {
	return getStrOption(name, globalParams)
}

// Returns the local string option's value.
// Returns an error in case the parameter name was incorrect.
func GetStrOption(name string) (string, error) {
	return getStrOption(name, localParams)
}

// Implements common functionality for obtaining an integer option value.
func getIntOption(name string, params map[string]ParamInfo) (int, error) {
	info, ok := params[name]
	if !ok {
		errorMsg := fmt.Sprintf("Unknown parameter name '%s'.", name)
		return 0, errors.New(errorMsg)
	}

	if info.kind != ParamInt {
		errorMsg := fmt.Sprintf("Parameter '%s' is not an integer option.", name)
		return 0, errors.New(errorMsg)
	}

	return info.valueInt, nil
}

// Returns the global integer option's value.
// Returns an error in case the parameter name was incorrect.
func GetGlobalIntOption(name string) (int, error) {
	return getIntOption(name, globalParams)
}

// Returns the local string option's value.
// Returns an error in case the parameter name was incorrect.
func GetIntOption(name string) (int, error) {
	return getIntOption(name, localParams)
}
