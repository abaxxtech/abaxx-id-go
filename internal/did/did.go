package did

import (
	"errors"
	"fmt"
	"regexp"
)

type DwnError struct {
	Code    string
	Message string
}

func (e *DwnError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

const (
	DidNotString = "DidNotString"
	DidNotValid  = "DidNotValid"
)

var didRegex = regexp.MustCompile(`^did:([a-z0-9]+):((?:(?:[a-zA-Z0-9._-]|(?:%[0-9a-fA-F]{2}))*:)*((?:[a-zA-Z0-9._-]|(?:%[0-9a-fA-F]{2}))+))((;[a-zA-Z0-9_.:%-]+=[a-zA-Z0-9_.:%-]*)*)(\/[^#?]*)?([?][^#]*)?(#.*)?$`)

type Did struct{}

// GetMethodSpecificId gets the method specific ID segment of a DID. ie. did:<method-name>:<method-specific-id>
func (d *Did) GetMethodSpecificId(did string) (string, error) {
	secondColonIndex := indexOfNth(did, ':', 2)
	if secondColonIndex == -1 {
		return "", errors.New("invalid DID format")
	}
	methodSpecificId := did[secondColonIndex+1:]
	return methodSpecificId, nil
}

// Validate validates the given DID
func (d *Did) Validate(did interface{}) error {
	didStr, ok := did.(string)
	if !ok {
		return &DwnError{Code: DidNotString, Message: fmt.Sprintf("DID is not string: %v", did)}
	}

	if !didRegex.MatchString(didStr) {
		return &DwnError{Code: DidNotValid, Message: fmt.Sprintf("DID is not a valid DID: %s", didStr)}
	}

	return nil
}

// GetMethodName gets the method name from a DID. ie. did:<method-name>:<method-specific-id>
func (d *Did) GetMethodName(did string) (string, error) {
	secondColonIndex := indexOfNth(did, ':', 2)
	if secondColonIndex == -1 {
		return "", errors.New("invalid DID format")
	}
	methodName := did[4:secondColonIndex]
	return methodName, nil
}

// indexOfNth finds the index of the nth occurrence of a character in a string
func indexOfNth(s string, char rune, n int) int {
	count := 0
	for i, c := range s {
		if c == char {
			count++
			if count == n {
				return i
			}
		}
	}
	return -1
}
