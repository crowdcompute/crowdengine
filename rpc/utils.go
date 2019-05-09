package rpc

import (
	"encoding/json"
	"errors"
	"strings"
)

// FindNamespaceMethod detect which function in which namespace is called
func FindNamespaceMethod(incomingMsg json.RawMessage) (string, string, error) {
	var in jsonRequest
	if err := json.Unmarshal(incomingMsg, &in); err != nil {
		return "", "", err
	}

	elems := strings.Split(in.Method, serviceMethodSeparator)
	if len(elems) != 2 {
		return "", "", errors.New("Method not given")
	}

	return elems[0], elems[1], nil
}
