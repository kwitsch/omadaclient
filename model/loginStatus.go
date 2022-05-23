package model

type LoginStatus struct {
	Login bool `json:"login"`
}

type LoginStatusResponse struct {
	ApiHeader
	Result LoginStatus `json:"result"`
}
