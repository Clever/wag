package validation

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type GlideTest struct {
	Title string
	Input string
	Error error
}

var glideYAMLTests = []GlideTest{
	{
		Title: "Success glide.yaml is up to date",
		Input: `import:
- package: github.com/lightstep/lightstep-tracer-go
  version: 0d48cd619841b1e1a3cdd20cd6ac97774c0002ce
- package: github.com/opentracing/opentracing-go
  version: ^1.0.0
- package: github.com/opentracing/basictracer-go
  version: 1b32af207119a14b1b231d451df3ed04a72efebf
- package: github.com/gorilla/mux
  version: 757bef944d0f21880861c2dd9c871ca543023cba
- package: github.com/golang/mock
  version: 13f360950a79f5864a972c786a10a50e44b69541
`,
		Error: nil,
	},
	{
		Title: "Error if glide.yaml out of date for one or more deps (lightstep-tracer-go)",
		Input: `import:
- package: github.com/lightstep/lightstep-tracer-go
  version: incorrect
- package: github.com/opentracing/opentracing-go
  version: ^1.0.0
- package: github.com/opentracing/basictracer-go
  version: 1b32af207119a14b1b231d451df3ed04a72efebf
- package: github.com/gorilla/mux
  version: 757bef944d0f21880861c2dd9c871ca543023cba
- package: github.com/golang/mock
  version: 13f360950a79f5864a972c786a10a50e44b69541
`,
		Error: &ListOfPeerDependencyError{
			Errors: []*PeerDependencyError{
				&PeerDependencyError{Package: "github.com/lightstep/lightstep-tracer-go", Version: "0d48cd619841b1e1a3cdd20cd6ac97774c0002ce", File: "glide.yaml"},
			},
		},
	},
	{
		Title: "Error if an item is missing altogether from glide.yaml (opentracing-go)",
		Input: `import:
- package: github.com/lightstep/lightstep-tracer-go
  version: 0d48cd619841b1e1a3cdd20cd6ac97774c0002ce
- package: github.com/opentracing/basictracer-go
  version: 1b32af207119a14b1b231d451df3ed04a72efebf
- package: github.com/gorilla/mux
  version: 757bef944d0f21880861c2dd9c871ca543023cba
- package: github.com/golang/mock
  version: 13f360950a79f5864a972c786a10a50e44b69541
`,
		Error: &ListOfPeerDependencyError{
			Errors: []*PeerDependencyError{
				&PeerDependencyError{Package: "github.com/opentracing/opentracing-go", Version: "^1.0.0", File: "glide.yaml"},
			},
		},
	},
}

var glideLockTests = []GlideTest{
	{
		Title: "Success if glide.lock is fully up to date",
		Input: `hash: fakehashfakehash2ab3a7fc19967467d0350d521f2efc13979838813aabe77b
updated: 2017-10-18T18:42:39.792248183Z
imports:
- name: github.com/lightstep/lightstep-tracer-go
  version: 0d48cd619841b1e1a3cdd20cd6ac97774c0002ce
- name: github.com/opentracing/opentracing-go
  version: ^1.0.0
- name: github.com/opentracing/basictracer-go
  version: 1b32af207119a14b1b231d451df3ed04a72efebf
- name: github.com/gorilla/mux
  version: 757bef944d0f21880861c2dd9c871ca543023cba
- name: github.com/golang/mock
  version: 13f360950a79f5864a972c786a10a50e44b69541
`,
		Error: nil,
	},
	{
		Title: "glide.lock up-to-date does not exact match semver versions, since it can't easily determine if semver == commit hash",
		Input: `hash: fakehashfakehash2ab3a7fc19967467d0350d521f2efc13979838813aabe77b
updated: 2017-10-18T18:42:39.792248183Z
imports:
- name: github.com/lightstep/lightstep-tracer-go
  version: 0d48cd619841b1e1a3cdd20cd6ac97774c0002ce
- name: github.com/opentracing/opentracing-go
  version: some-unknown-version
- name: github.com/opentracing/basictracer-go
  version: 1b32af207119a14b1b231d451df3ed04a72efebf
- name: github.com/gorilla/mux
  version: 757bef944d0f21880861c2dd9c871ca543023cba
- name: github.com/golang/mock
  version: 13f360950a79f5864a972c786a10a50e44b69541
`,
		Error: nil,
	},
	{
		Title: "Error if glide.lock out of date for one or more deps (lightstep-tracer-go)",
		Input: `hash: fakehashfakehash2ab3a7fc19967467d0350d521f2efc13979838813aabe77b
updated: 2017-10-18T18:42:39.792248183Z
imports:
- name: github.com/lightstep/lightstep-tracer-go
  version: incorrect
- name: github.com/opentracing/opentracing-go
  version: ^1.0.0
- name: github.com/opentracing/basictracer-go
  version: 1b32af207119a14b1b231d451df3ed04a72efebf
- name: github.com/gorilla/mux
  version: 757bef944d0f21880861c2dd9c871ca543023cba
- name: github.com/golang/mock
  version: 13f360950a79f5864a972c786a10a50e44b69541
`,
		Error: &ListOfPeerDependencyError{
			Errors: []*PeerDependencyError{
				&PeerDependencyError{Package: "github.com/lightstep/lightstep-tracer-go", Version: "0d48cd619841b1e1a3cdd20cd6ac97774c0002ce", File: "glide.lock"},
			},
		},
	},
	{
		Title: "Error if an item is missing altogether from glide.lock (opentracing-go)",
		Input: `hash: fakehashfakehash2ab3a7fc19967467d0350d521f2efc13979838813aabe77b
updated: 2017-10-18T18:42:39.792248183Z
imports:
- name: github.com/lightstep/lightstep-tracer-go
  version: 0d48cd619841b1e1a3cdd20cd6ac97774c0002ce
- name: github.com/opentracing/basictracer-go
  version: 1b32af207119a14b1b231d451df3ed04a72efebf
- name: github.com/gorilla/mux
  version: 757bef944d0f21880861c2dd9c871ca543023cba
- name: github.com/golang/mock
  version: 13f360950a79f5864a972c786a10a50e44b69541
`,
		Error: &ListOfPeerDependencyError{
			Errors: []*PeerDependencyError{
				&PeerDependencyError{Package: "github.com/opentracing/opentracing-go", Version: "^1.0.0", File: "glide.lock"},
			},
		},
	},
}

func TestValidateGlideYAML(t *testing.T) {
	for _, test := range glideYAMLTests {
		t.Run(test.Title, func(t *testing.T) {
			err := ValidateGlideYAML(strings.NewReader(test.Input))
			if test.Error != nil {
				assert.Error(t, err)
				assert.Equal(t, test.Error, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateGlideLock(t *testing.T) {
	for _, test := range glideLockTests {
		t.Log(test.Title)
		err := ValidateGlideLock(strings.NewReader(test.Input))
		if test.Error != nil {
			assert.Error(t, err)
			assert.Equal(t, test.Error, err)
		} else {
			assert.NoError(t, err)
		}
	}
}
