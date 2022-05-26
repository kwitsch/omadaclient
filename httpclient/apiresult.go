package httpclient

import (
	"encoding/json"
)

type ApiResult struct {
	ErrorCode int         `json:"errorCode"`
	Msg       string      `json:"msg"`
	Result    interface{} `json:"result"`
}

func (ar ApiResult) GetHead() (int, string) {
	return ar.ErrorCode, ar.Msg
}

func (ar ApiResult) IsSuccess() bool {
	return (ar.ErrorCode == 0)
}

func (ar ApiResult) GetResult(res interface{}) error {
	b, err := json.Marshal(&ar.Result)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, &res)
}
