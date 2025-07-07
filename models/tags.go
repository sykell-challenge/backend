package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Tag struct {
	TagName string `json:"tag_name"`
	Count   int    `json:"count"`
}

type Tags []Tag

// Value implements the driver.Valuer interface
func (t Tags) Value() (driver.Value, error) {
	return json.Marshal(t)
}

// Scan implements the sql.Scanner interface
func (t *Tags) Scan(value interface{}) error {
	if value == nil {
		*t = Tags{}
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, t)
	case string:
		return json.Unmarshal([]byte(v), t)
	default:
		return fmt.Errorf("cannot scan %T into Tags", value)
	}
}
