package model

type ApiInfo struct {
	ControllerVer string `json:"controllerVer"`
	ApiVer        string `json:"apiVer"`
	Configured    bool   `json:"configured"`
	Type          int    `json:"type"`
	SupportApp    bool   `json:"supportApp"`
	OmadacId      string `json:"omadacId"`
}

type ApiInfoResponse struct {
	ErrorCode int     `json:"errorCode"`
	Msg       string  `json:"msg"`
	Result    ApiInfo `json:"result"`
}
