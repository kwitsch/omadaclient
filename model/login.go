package model

type Login struct {
	RoleType int    `json:"roleType"`
	Token    string `json:"token"`
}

type LoginResponse struct {
	ApiResponse
	Result Login `json:"result"`
}
