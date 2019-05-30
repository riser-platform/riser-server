package core

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// jsonbSqlUnmarshal unmarshals jsonb from Postgres. Since calling code conforms to the SQL scan interface, the value coming in
// must be interface{} even though we know if it's jsonb it will always be a []byte
func jsonbSqlUnmarshal(v interface{}, t interface{}) error {
	if v == nil {
		return nil
	}
	b, ok := v.([]byte)
	if !ok {
		return errors.New("unable to read jsonb field as []byte")
	}

	return json.Unmarshal(b, &t)
}
