package util

import (
	"time"

	"github.com/google/uuid"
)

func FormatDateTime(ts time.Time) string {
	return ts.Format(time.DateTime)
}

func UUID() string {
	return uuid.New().String()
}
