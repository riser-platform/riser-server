package model

import (
	validation "github.com/go-ozzo/ozzo-validation/v3"
	"github.com/google/uuid"
)

type AppName string

type App struct {
	Id        uuid.UUID     `json:"id"`
	Name      AppName       `json:"name"`
	Namespace NamespaceName `json:"namespace"`
}

func (v App) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.Name),
		validation.Field(&v.Namespace))
}

type NewApp struct {
	Name      AppName       `json:"name"`
	Namespace NamespaceName `json:"namespace"`
}

func (v NewApp) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.Name),
		validation.Field(&v.Namespace))
}

func (v AppName) Validate() error {
	return validation.Validate(string(v), RulesAppName()...)
}
