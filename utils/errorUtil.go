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

	errMsg := fmt.Sprintf("%s:\n %s",
		baseErr.Error(), strings.Join(errStrings, " :: "))
	return errors.New(errMsg)
}
