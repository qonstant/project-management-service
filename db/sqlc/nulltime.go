package db

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type NullableTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface.
func (nt *NullableTime) Scan(value interface{}) error {
	if value == nil {
		nt.Valid = false
		return nil
	}
	nt.Valid = true
	return convertTime(value, &nt.Time)
}

// Value implements the Valuer interface.
func (nt NullableTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

// convertTime converts an interface{} to time.Time
func convertTime(value interface{}, t *time.Time) error {
	switch v := value.(type) {
	case time.Time:
		*t = v
		return nil
	case nil:
		*t = time.Time{}
		return nil
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
}
