package model

type ApiInfo struct {
	ControllerVer string `json:"controllerVer"`
	ApiVer        string `json:"apiVer"`
	Configured    bool   `json:"configured"`
	Type          uint   `json:"type"`
	SupportApp    bool   `json:"supportApp"`
	OmadacId      string `json:"omadacId"`
}
