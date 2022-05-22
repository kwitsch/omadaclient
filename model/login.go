package model

type Login struct {
	RoleType int    `json:"roleType"`
	Token    string `json:"token"`
}

type LoginResponse struct {
	ErrorCode int    `json:"errorCode"`
	Msg       string `json:"msg"`
	Result    Login  `json:"result"`
}
