package model

import "github.com/google/uuid"

type UnsealedSecret struct {
	SecretMeta `json:",inline"`
	PlainText  string `json:"secretValue"`
}

type SecretMeta struct {
	AppId uuid.UUID `json:"appId"`
	Stage string    `json:"stage"`
	Name  string    `json:"secretName"`
}

type SecretMetaStatus struct {
	SecretMeta `json:",inline"`
	Revision   int64 `json:"revision"`
}
