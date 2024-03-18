package helpers

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Tags []string

func (t *Tags) Scan(value interface{}) error {
	if value == nil {
		*t = make(Tags, 0)
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, t)
	case string:
		return json.Unmarshal([]byte(v), t)
	default:
		return errors.New("unsupported type")
	}
}

func (t Tags) Value() (driver.Value, error) {
	return json.Marshal(t)
}
