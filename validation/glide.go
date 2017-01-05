package validation

import (
	"fmt"
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// GlideYML unmarshals the parts of a glide.yml file we care about.
type GlideYML struct {
	Imports []Import `yaml:"import"`
}

// Import contained within a glide.yml.
type Import struct {
	Package string `yaml:"package"`
	Version string `yaml:"version"`
}

// requirements describes dependencies we require wag users to use in their apps.
var requirements = []Import{
	{
		Package: "github.com/lightstep/lightstep-tracer-go",
		Version: "0d48cd619841b1e1a3cdd20cd6ac97774c0002ce",
	},
	{
		Package: "github.com/opentracing/opentracing-go",
		Version: "^1.0.0",
	},
	{
		Package: "github.com/opentracing/basictracer-go",
		Version: "1b32af207119a14b1b231d451df3ed04a72efebf",
	},
	{
		Package: "github.com/gorilla/mux",
		Version: "757bef944d0f21880861c2dd9c871ca543023cba",
	},
}

// ValidateGlideYML looks at a user's glide.yml and makes sure certain dependencies
// that wag requires are present and locked to the correct version.
func ValidateGlideYML(glideYMLFile io.Reader) error {
	var glideYML GlideYML
	bs, err := ioutil.ReadAll(glideYMLFile)
	if err != nil {
		return fmt.Errorf("error reading glide.yml: %s", err)
	}
	if err = yaml.Unmarshal(bs, &glideYML); err != nil {
		return fmt.Errorf("error unmarshalling yaml: %s", err)
	}

	for _, req := range requirements {
		if err := validateImports(glideYML.Imports, req); err != nil {
			return err
		}
	}

	return nil

}

func validateImports(imports []Import, requiredImport Import) error {
	for _, i := range imports {
		if i.Package == requiredImport.Package &&
			i.Version == requiredImport.Version {
			return nil
		}
	}
	return fmt.Errorf("wag requires version %s of %s. Please update your glide.yml and run `glide up`", requiredImport.Version, requiredImport.Package)
}
