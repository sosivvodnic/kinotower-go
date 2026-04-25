package domain

import (
	"time"

	"github.com/sosivvodnic/kinotower-go/internal/core/httpjson"
)

func NewISOTime(t time.Time) httpjson.ISOTime {
	return httpjson.ISOTime{Time: t}
}

