package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Link struct {
	Link       string `json:"link"`
	Type       string `json:"type"` // internal, external, inaccessible
	StatusCode int    `json:"status_code"`
}

type Links []Link

// Value implements the driver.Valuer interface
func (l Links) Value() (driver.Value, error) {
	return json.Marshal(l)
}

// Scan implements the sql.Scanner interface
func (l *Links) Scan(value interface{}) error {
	if value == nil {
		*l = Links{}
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, l)
	case string:
		return json.Unmarshal([]byte(v), l)
	default:
		return fmt.Errorf("cannot scan %T into Links", value)
	}
}
