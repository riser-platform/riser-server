package model

type App struct {
	// Id is a short hash (32bit SHA1 represented as an 8 character hex)
	Id   string `json:"id"`
	Name string `json:"name"`
}

type NewApp struct {
	Name string `json:"name"`
}
