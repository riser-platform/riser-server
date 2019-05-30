package core

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type Stage struct {
	Name string
	Doc  StageDoc
}

type StageDoc struct {
	LastPing time.Time   `json:"lastPing"`
	Config   StageConfig `json:"config"`
}

type StageConfig struct {
	SealedSecretCert  []byte `json:"sealedSecretCert"`
	PublicGatewayHost string `json:"publicGatewayHost"`
}

// Needed for sql.Scanner interface
func (a *StageDoc) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Needed for sql.Scanner interface
func (a *StageDoc) Scan(value interface{}) error {
	return jsonbSqlUnmarshal(value, &a)
}
