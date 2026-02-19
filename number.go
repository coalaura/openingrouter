package openingrouter

import (
	"bytes"
	"encoding/json"

	"github.com/coalaura/byteconv"
)

// StringifiedNumber represents a number that was stringified in JSON
type StringifiedNumber float64

func (sn *StringifiedNumber) UnmarshalJSON(data []byte) error {
	// Handle null as zero
	if bytes.Equal(data, []byte("null")) || bytes.Equal(data, []byte("\"\"")) {
		*sn = 0

		return nil
	}

	if len(data) > 0 && data[0] == '"' {
		f, err := byteconv.ParseFloat(data[1:len(data)-1], 64)
		if err != nil {
			return err
		}

		*sn = StringifiedNumber(f)

		return nil
	}

	// Try as a plain number
	f, err := byteconv.ParseFloat(data, 64)
	if err != nil {
		return err
	}

	*sn = StringifiedNumber(f)

	return nil
}

func (sn StringifiedNumber) MarshalJSON() ([]byte, error) {
	return json.Marshal(float64(sn))
}

// Float64 returns the float64 value
func (sn StringifiedNumber) Float64() float64 {
	return float64(sn)
}
