package model

type ApiHeader struct {
	ApiHeaderMethodes
	ErrorCode int    `json:"errorCode"`
	Msg       string `json:"msg"`
}

type ApiHeaderMethodes interface {
	GetHead() *ApiHeader
	IsSuccess() bool
}

func (ah ApiHeader) GetHead() *ApiHeader {
	return &ah
}

func (ah ApiHeader) IsSuccess() bool {
	return (ah.ErrorCode == 0)
}
