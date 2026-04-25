package httpjson

import (
	"time"
)

// ISOTime formats timestamps as ISO-8601 with milliseconds and Z suffix:
// 2032-01-31T21:59:35.000Z
type ISOTime struct {
	time.Time
}

func (t ISOTime) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte("null"), nil
	}
	s := t.Time.UTC().Format("2006-01-02T15:04:05.000Z")
	return []byte(`"` + s + `"`), nil
}

