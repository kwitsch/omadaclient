package model

type ApiResponse struct {
	ErrorCode int    `json:"errorCode"`
	Msg       string `json:"msg"`
}
