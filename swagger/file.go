package swagger

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/go-openapi/spec"
)

// WriteToFile writes the swagger spec to a temporary file. It returns
// the name of the file.
func WriteToFile(s *spec.Swagger) (string, error) {
	bytes, err := json.Marshal(s)
	if err != nil {
		return "", fmt.Errorf("internal error: wag patch type marshal failure: %s", err)
	}
	// TODO: Better way to do temporary files... I think there's one in Go
	fileName := "swagger.tmp"
	if err := ioutil.WriteFile(fileName, bytes, 0644); err != nil {
		return "", err
	}
	return fileName, nil
}
