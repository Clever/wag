package swagger

import (
	"encoding/json"
	"log"

	"github.com/go-openapi/spec"
)

const SubrouterKey string = "x-routers"

type Subrouter struct {
	Key  string `json:"key"`
	Path string `json:"path"`
}

func ParseSubrouters(s spec.Swagger) ([]Subrouter, error) {
	var subrouterConfig []Subrouter
	if routers, ok := s.Extensions[SubrouterKey]; ok {
		if subroutersM, ok := routers.([]interface{}); ok {
			subroutersB, err := json.Marshal(subroutersM)
			if err != nil {
				return nil, err
			}

			err = json.Unmarshal(subroutersB, &subrouterConfig)
			if err != nil {
				return nil, err
			}
		} else {
			log.Printf("WARNING: %s subrouter config was not an array\n", SubrouterKey)
		}
	}

	return subrouterConfig, nil
}
