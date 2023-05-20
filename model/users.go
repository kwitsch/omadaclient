package model

type UserPrivilege struct {
	All         bool            `json:"all"`
	LastVisited string          `json:"lastVisited"`
	Sites       []PrivilegeSite `json:"sites"`
}

type PrivilegeSite struct {
	Key     string `json:"key"`
	Name    string `json:"name"`
	Primary bool   `json:"primary"`
}

type UsersCurrent struct {
	Id        string        `json:"id"`
	Type      uint64        `json:"type"`
	RoleType  uint64        `json:"roleType"`
	Privilege UserPrivilege `json:"privilege"`
	Adopt     bool          `json:"adopt"`
	Manage    bool          `json:"manage"`
	License   bool          `json:"license"`
	Name      string        `json:"name"`
	Email     string        `json:"email"`
	Alert     bool          `json:"alert"`
	Disaster  int           `json:"disaster"`
	Favorites []string      `json:"favorites"`
	Dbnormal  bool          `json:"dbnormal"`
}
