package core

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type Environment struct {
	Name string
	Doc  EnvironmentDoc
}

type EnvironmentDoc struct {
	LastPing time.Time         `json:"lastPing"`
	Config   EnvironmentConfig `json:"config"`
}

type EnvironmentConfig struct {
	SealedSecretCert  []byte `json:"sealedSecretCert"`
	PublicGatewayHost string `json:"publicGatewayHost"`
}

// Needed for sql.Scanner interface
func (a *EnvironmentDoc) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Needed for sql.Scanner interface
func (a *EnvironmentDoc) Scan(value interface{}) error {
	return jsonbSqlUnmarshal(value, &a)
}

type EnvironmentStatus struct {
	EnvironmentName string
	Healthy         bool
	Reason          string
}
