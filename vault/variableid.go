package vault

import (
	"fmt"
	"strings"
)

type VariableID struct {
	Path string
	Key  string
}

func NewVariableID(argument string) (*VariableID, error) {
	var path, key, err = sanitize(argument)

	if err != nil {
		return nil, err
	}
	valuePath := &VariableID{
		Path: path,
		Key:  key,
	}
	return valuePath, nil
}

func sanitize(argument string) (string, string, error) {
	argument = strings.TrimSpace(argument)
	countCross := strings.Count(argument, "#")
	if countCross > 1 {
		return "", "", fmt.Errorf("SYNTAX ERROR: variableID %q contains %d '#'. Maximum ONE '#' is allowed", argument, countCross)
	}

	var path = argument
	var key string
	if countCross == 1 {
		var arguments = strings.SplitN(argument, "#", 2)
		path = arguments[0]
		key = arguments[1]
		if strings.Count(key, "/") > 0 {
			return "", "", fmt.Errorf("SYNTAX ERROR: variableID %q (AWS format) with a slash after the '#'", argument)
		}
		if len(key) < 1 {
			return "", "", fmt.Errorf("SYNTAX ERROR: variableID %q (AWS format) ends with '#'", argument)
		}
	}

	var parts = strings.Split(path, "/")
	if len(parts) == 1 {
		return "", "", fmt.Errorf("SYNTAX ERROR: variableID %q DOESN'T contain any slash or '#'", argument)
	}

	if len(key) == 0 {
		key = parts[len(parts)-1]
		// remove key form parts
		parts = parts[:len(parts)-1]
	}

	if len(key) < 1 {
		return "", "", fmt.Errorf("SYNTAX ERROR: variableID %q key is empty", argument)
	}

	for i := 0; i < len(parts); i++ {
		if len(parts[i]) < 1 {
			return "", "", fmt.Errorf("SYNTAX ERROR: variableID %q contains leading slash or ending slash or double slashes", argument)
		}
	}

	return strings.Join(parts, "/"), key, nil
}
