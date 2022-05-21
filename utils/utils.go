package utils

import (
	"errors"
	"fmt"
)

func NewError(input ...interface{}) error {
	return errors.New(fmt.Sprint(input...))
}
