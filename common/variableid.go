package common

import (
	"fmt"
	"strings"
)

type VariableID struct {
	Path string
	Key  string
}

func NewVariableID(argument string) (*VariableID, error) {
	var path, key, err = Sanitize(argument)

	if err != nil {
		return nil, err
	}
	valuePath := &VariableID{
		Path: path,
		Key:  key,
	}
	return valuePath, nil
}

func Sanitize(argument string) (string, string, error) {
	argument = strings.TrimSpace(argument)
	countCross := strings.Count(argument, "#")
	if countCross > 1 {
		return "", "", fmt.Errorf("variableID %q contains %d '#'. Maximum ONE '#' is allowed", argument, countCross)
	}

	var path = argument
	var key string
	if countCross == 1 {
		var arguments = strings.SplitN(argument, "#", 2)
		path = arguments[0]
		key = arguments[1]
		if strings.Count(key, "/") > 0 {
			return "", "", fmt.Errorf("variableID %q AWS format with a '/' after the '#'", argument)
		}
		if len(key) < 1 {
			return "", "", fmt.Errorf("variableID %q AWS format ends with '#'", argument)
		}

	}

	// sanitize for leading, trailing and concatenated '/'
	var parts = strings.Split(path, "/")

	var length = len(parts)
	if length == 1 {
		return "", "", fmt.Errorf("variableID %q DOESN'T contain any '/' or '#'", argument)
	}
	path = ""

	if len(key) == 0 {
		length--
		key = parts[length]
	}

	for i := 0; i < length; i++ {
		if len(parts[i]) < 1 {
			return "", "", fmt.Errorf("argument %q contains leading slash or ending slash or double //", argument)
		}
		path += parts[i]
		if i != length-1 {
			path += "/"
		}
	}

	if len(key) < 1 {
		return "", "", fmt.Errorf("variableID %q key is empty", argument)
	}
	return path, key, nil
}
