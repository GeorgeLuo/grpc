package main

import (
	"errors"
	"fmt"
)

// AppendError returns a trace error chain with a new error appended
func AppendError(baseErr error, newErr error) error {
	errMsg := fmt.Sprintf("%s:\n %s",
		baseErr.Error(), newErr.Error())
	return errors.New(errMsg)
}
