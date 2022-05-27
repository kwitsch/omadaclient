package model

type UsersCurrent struct {
	Id        string `json:"id"`
	Type      int    `json:"type"`
	RoleType  int    `json:"roleType"`
	Privilege struct {
		All         bool   `json:"all"`
		LastVisited string `json:"lastVisited"`
		Sites       []struct {
			Key     string `json:"key"`
			Name    string `json:"name"`
			Primary bool   `json:"primary"`
		} `json:"sites"`
	} `json:"privilege"`
	Adopt     bool     `json:"adopt"`
	Manage    bool     `json:"manage"`
	License   bool     `json:"license"`
	Name      string   `json:"name"`
	Email     string   `json:"email"`
	Alert     bool     `json:"alert"`
	Disaster  int      `json:"disaster"`
	Favorites []string `json:"favorites"`
	Dbnormal  bool     `json:"dbnormal"`
}