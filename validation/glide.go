package validation

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"

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

// GlideLock unmarshals the parts of a glide.lock file we care about.
type GlideLock struct {
	Imports []LockedVersion `yaml:"imports"`
}

// LockedVersion contained within a glide.lock.
type LockedVersion struct {
	Name    string `yaml:"name"`
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
	{
		Package: "github.com/golang/mock",
		Version: "13f360950a79f5864a972c786a10a50e44b69541",
	},
}

// ValidateGlideYAML looks at a user's glide.yml and makes sure certain dependencies
// that wag requires are present and locked to the correct version.
func ValidateGlideYAML(glideYMLFile io.Reader) error {
	var glideYML GlideYML
	bs, err := ioutil.ReadAll(glideYMLFile)
	if err != nil {
		return fmt.Errorf("error reading glide.yaml: %s", err)
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

// ValidateGlideLock looks at a user's glide.yml and makes sure certain dependencies
// that wag requires are present and locked to the correct version.
func ValidateGlideLock(glideLockFile io.Reader) error {
	var glideLock GlideLock
	bs, err := ioutil.ReadAll(glideLockFile)
	if err != nil {
		return fmt.Errorf("error reading glide.lock: %s", err)
	}
	if err = yaml.Unmarshal(bs, &glideLock); err != nil {
		return fmt.Errorf("error unmarshalling yaml in glide.lock: %s", err)
	}

	for _, req := range requirements {
		if err := validateLockedVersion(glideLock.Imports, req); err != nil {
			return err
		}
	}

	return nil
}

// PeerDependencyError occurs when glide.yml and/or glide.lock dont have the
// required dependency versions for wag
type PeerDependencyError struct {
	Package string
	Version string
	File    string
}

func (e *PeerDependencyError) Error() string {
	return fmt.Sprintf("Error: wag peer dependency not met. \n"+
		"Version %s of %s must be set in glide.yaml and glide.lock.\n"+
		"Please ensure your glide.yaml file includes\n\n"+
		"```\n"+
		"- package: %s\n"+
		"  version: %s\n"+
		"```\n\n"+
		"then run `glide up`.", e.Version, e.Package, e.Package, e.Version)
}

func validateImports(imports []Import, requiredImport Import) error {
	for _, i := range imports {
		if i.Package == requiredImport.Package &&
			i.Version == requiredImport.Version {
			return nil
		}
	}
	return &PeerDependencyError{Package: requiredImport.Package, Version: requiredImport.Version, File: "glide.yaml"}
}

func validateLockedVersion(versions []LockedVersion, requiredImport Import) error {
	for _, v := range versions {
		// If we've specified locking to a semantic version like "^1.0.0", we can't easily determine
		// if the locked version satisfies that. We accept any version as long as the package is present.
		if (v.Name == requiredImport.Package && strings.HasPrefix(requiredImport.Version, "^")) ||
			(v.Name == requiredImport.Package && v.Version == requiredImport.Version) {
			return nil
		}
	}
	return &PeerDependencyError{Package: requiredImport.Package, Version: requiredImport.Version, File: "glide.lock"}
}
