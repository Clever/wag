package swagger

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-openapi/spec"
)

// WriteToFile writes the swagger spec to a temporary file. It returns
// the name of the created file.
func WriteToFile(s *spec.Swagger) (string, error) {
	bytes, err := json.Marshal(s)
	if err != nil {
		return "", fmt.Errorf("internal error: wag patch type marshal failure: %s", err)
	}
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		return "", err
	}

	if _, err := tmpfile.Write(bytes); err != nil {
		return "", err
	}
	return tmpfile.Name(), nil
}
