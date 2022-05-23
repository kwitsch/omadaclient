package model

type LoginStatus struct {
	Login bool `json:"login"`
}

type LoginStatusResponse struct {
	ApiResponse
	Result LoginStatus `json:"result"`
}
