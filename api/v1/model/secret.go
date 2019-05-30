package model

import (
	"time"
)

type UnsealedSecret struct {
	SecretMeta `json:",inline"`
	PlainText  string `json:"secretValue"`
}

type SecretMeta struct {
	App   string `json:"app"`
	Stage string `json:"stage"`
	Name  string `json:"secretName"`
}

type SecretMetaStatus struct {
	SecretMeta  `json:",inline"`
	LastUpdated time.Time `json:"lastUpdated"`
}
