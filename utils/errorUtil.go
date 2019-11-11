package utils

import (
	"errors"
	"fmt"
	"strings"
)

// AppendError returns a trace error chain with a new error appended
func AppendError(baseErr error, newErr ...error) error {
	var errStrings []string
	for _, err := range newErr {
		errStrings = append(errStrings, err.Error())
	}

	errMsg := fmt.Sprintf("%s\n%s",
		baseErr.Error(), strings.Join(errStrings, "\n"))
	return errors.New(errMsg)
}

// AppendStringToError returns a trace error chain with a new error as a
// string appended
func AppendStringToError(baseErr error, newErr ...string) error {
	var errStrings []string
	for _, err := range newErr {
		errStrings = append(errStrings, err)
	}

	errMsg := fmt.Sprintf("%s\n%s",
		baseErr.Error(), strings.Join(errStrings, "\n"))
	return errors.New(errMsg)
}
