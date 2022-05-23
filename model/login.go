package model

type Login struct {
	RoleType int    `json:"roleType"`
	Token    string `json:"token"`
}

type LoginResponse struct {
	ApiHeader
	Result Login `json:"result"`
}
