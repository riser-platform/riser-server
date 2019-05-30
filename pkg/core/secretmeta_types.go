package core

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type SecretMeta struct {
	AppName    string
	StageName  string
	SecretName string
	Doc        SecretMetaDoc
}

type SecretMetaDoc struct {
	LastUpdated time.Time `json:"lastUpdated"`
}

// Needed for sql.Scanner interface
func (a *SecretMetaDoc) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Needed for sql.Scanner interface
func (a *SecretMetaDoc) Scan(value interface{}) error {
	return jsonbSqlUnmarshal(value, &a)
}
