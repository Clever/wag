package swagger

import (
	"encoding/hex"

	"encoding/hex"
)

// InitCustomFormats adds wag's custom formats to the global go-openapi/strfmt Default registry.
func InitCustomFormats() {
	m := mongoID("")
	strfmt.Default.Add("mongo-id", &m, isMongoID)
}

func isMongoID(s string) bool {
	if len(s) != 24 {
		return false
	}
	_, err := hex.DecodeString(s)
	return err == nil
}

type mongoID string

// MarshalText turns this instance into text
func (e mongoID) MarshalText() ([]byte, error) {
	return []byte(string(e)), nil
}

// UnmarshalText hydrates this instance from text
func (e *mongoID) UnmarshalText(data []byte) error { // validation is performed later on
	*e = mongoID(string(data))
	return nil
}

// String representation of the Mongo ID.
func (e mongoID) String() string {
	return string(e)
}
