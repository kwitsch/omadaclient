package model

type UsersCurrent struct {
	Id        string                `json:"id"`
	Type      int                   `json:"type"`
	RoleType  int                   `json:"roleType"`
	Privilege UsersCurrentPrivilege `json:"privilege"`
}

type UsersCurrentPrivilege struct {
	All         bool                        `json:"all"`
	LastVisited string                      `json:"lastVisited"`
	Sites       []UsersCurrentPrivilegeSite `json:"sites"`
}

type UsersCurrentPrivilegeSite struct {
}
