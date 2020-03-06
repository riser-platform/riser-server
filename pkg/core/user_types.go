package core

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id       uuid.UUID
	Username string
	Doc      UserDoc
}

type UserDoc struct {
	Created time.Time `json:"created"`
}

type NewUser struct {
	Id       uuid.UUID
	Username string
}

// Needed for sql.Scanner interface
func (a *UserDoc) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Needed for sql.Scanner interface
func (a *UserDoc) Scan(value interface{}) error {
	return jsonbSqlUnmarshal(value, &a)
}
