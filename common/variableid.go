package common

import (
	"fmt"
	"strings"
)

type VariableID struct {
	Path string
	Key  string
}

func NewVariableID(argument string) *VariableID {
	var path, key = Sanitize(argument)

	valuePath := &VariableID{
		Path: path,
		Key:  key,
	}
	return valuePath
}

func Sanitize(s string) (string, string) {
	s = strings.TrimSpace(s)
	countCross := strings.Count(s, "#")
	if countCross > 1 {
		Exit(fmt.Errorf("variableID %q contains %d '#'. Maximum ONE '#' is allowed", s, countCross))
	}

	var path = s
	var key string
	if countCross == 1 {
		var arguments = strings.SplitN(s, "#", 2)
		path = arguments[0]
		key = arguments[1]
	}

	// sanitize for leading, trailing and concatenated '/'
	var parts = strings.Split(path, "/")

	var length = len(parts)
	if len(key) == 0 && length == 1 {
		Exit(fmt.Errorf("variableID %q DOESN'T contain any '/' or '#'", s))
	}
	path = ""

	if len(key) == 0 {
		length--
		key = parts[length]
	}

	for i := 0; i < length; i++ {
		path += parts[i]
		if i != length-1 {
			path += "/"
		}
	}

	return path, key
}
