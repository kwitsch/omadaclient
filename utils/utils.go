package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

func NewError(input ...interface{}) error {
	return errors.New(ToString(input...))
}

func ToString(params ...interface{}) string {
	result := ""
	for _, p := range params {
		if len(result) > 0 {
			result += " "
		}
		if b, ok := p.([]byte); ok {
			result += string(b)
		} else {
			rv := reflect.ValueOf(p)
			if rv.Kind() == reflect.Struct {
				if b, err := json.Marshal(p); err == nil {
					result += string(b)
				}
			} else {
				result += fmt.Sprint(p)
			}
		}
	}

	return result
}
