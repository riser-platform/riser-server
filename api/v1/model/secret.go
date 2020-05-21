package model

import validation "github.com/go-ozzo/ozzo-validation/v3"

type UnsealedSecret struct {
	SecretMeta `json:",inline"`
	PlainText  string `json:"plainTextValue"`
}

func (v UnsealedSecret) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.SecretMeta),
		validation.Field(&v.PlainText, validation.Required))
}

type SecretMeta struct {
	Name        string        `json:"name"`
	AppName     AppName       `json:"app"`
	Namespace   NamespaceName `json:"namespace"`
	Environment string        `json:"environment"`
}

func (v SecretMeta) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.AppName),
		validation.Field(&v.Namespace),
		validation.Field(&v.Name, validation.Required),
		validation.Field(&v.Environment, validation.Required))
}

type SecretMetaStatus struct {
	SecretMeta `json:",inline"`
	Revision   int64 `json:"revision"`
}
