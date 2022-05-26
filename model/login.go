package model

type Login struct {
	RoleType int    `json:"roleType"`
	Token    string `json:"token"`
}

type LoginStatus struct {
	Login bool `json:"login"`
}
