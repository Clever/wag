package swagger

import (
	"regexp"

	"github.com/go-openapi/strfmt"
)

func InitCustomFormats() {
	m := mongoID("")
	strfmt.Default.Add("mongoID", &m, isMongoID)
}

var mongoRegExp = regexp.MustCompile("^[0-9a-f]{24}$")

func isMongoID(str string) bool {
	return mongoRegExp.MatchString(str)
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

func (e mongoID) String() string {
	return string(e)
}
