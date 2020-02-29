package model

import (
	"github.com/google/uuid"
)

type App struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type NewApp struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}
