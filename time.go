package openingrouter

import (
	"encoding/json"
	"strings"
	"time"
)

// FlexibleTime handles multiple timestamp formats
type FlexibleTime struct {
	time.Time
}

var timeFormats = []string{
	"2006-01-02T15:04:05.999999999Z07:00", // RFC3339 with nanoseconds
	"2006-01-02T15:04:05Z07:00",           // RFC3339 without nanoseconds
	"2006-01-02T15:04:05Z",                // RFC3339 with Z
	"2006-01-02",                          // Date only
}

func (ft *FlexibleTime) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)
	if str == "null" || str == "" {
		ft.Time = time.Time{}

		return nil
	}

	var err error

	for _, format := range timeFormats {
		ft.Time, err = time.Parse(format, str)
		if err == nil {
			return nil
		}
	}

	return err
}

func (ft FlexibleTime) MarshalJSON() ([]byte, error) {
	if ft.Time.IsZero() {
		return []byte("null"), nil
	}

	return json.Marshal(ft.Time.Format(time.RFC3339))
}
