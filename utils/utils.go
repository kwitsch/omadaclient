package utils

import (
	"encoding/json"
	"errors"
	"fmt"
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
		switch v := p.(type) {
		case []byte:
			result += string(v)
		case string:
			result += v
		case int:
			result += fmt.Sprint(v)
		case uint:
			result += fmt.Sprint(v)
		case bool:
			result += fmt.Sprint(v)
		case float32:
			result += fmt.Sprint(v)
		case float64:
			result += fmt.Sprint(v)
		case error:
			result += v.Error()
		default:
			if b, err := json.Marshal(v); err == nil {
				result += string(b)
			}
		}
	}

	return result
}
