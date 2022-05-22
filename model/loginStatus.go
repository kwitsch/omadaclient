package model

type LoginStatus struct {
	Login bool `json:"login"`
}

type LoginStatusResponse struct {
	ErrorCode int         `json:"errorCode"`
	Msg       string      `json:"msg"`
	Result    LoginStatus `json:"result"`
}
