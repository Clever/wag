package validation

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type GlideYMLTest struct {
	YML   string
	Error error
}

var glideYMLTests = []GlideYMLTest{
	{
		YML: `import:
- package: github.com/lightstep/lightstep-tracer-go
  version: 0d48cd619841b1e1a3cdd20cd6ac97774c0002ce
- package: github.com/opentracing/opentracing-go
  version: ^1.0.0
- package: github.com/opentracing/basictracer-go
  version: 1b32af207119a14b1b231d451df3ed04a72efebf
- package: github.com/gorilla/mux
  version: 757bef944d0f21880861c2dd9c871ca543023cba
`,
		Error: nil,
	},
	{
		YML: `import:
- package: github.com/lightstep/lightstep-tracer-go
  version: incorrect
- package: github.com/opentracing/opentracing-go
  version: ^1.0.0
- package: github.com/opentracing/basictracer-go
  version: 1b32af207119a14b1b231d451df3ed04a72efebf
- package: github.com/gorilla/mux
  version: 757bef944d0f21880861c2dd9c871ca543023cba
`,
		Error: errors.New("wag requires version 0d48cd619841b1e1a3cdd20cd6ac97774c0002ce of github.com/lightstep/lightstep-tracer-go. Please update your glide.yml and run `glide up`"),
	},
	{
		YML: `import:
- package: github.com/lightstep/lightstep-tracer-go
  version: 0d48cd619841b1e1a3cdd20cd6ac97774c0002ce
- package: github.com/opentracing/basictracer-go
  version: 1b32af207119a14b1b231d451df3ed04a72efebf
- package: github.com/gorilla/mux
  version: 757bef944d0f21880861c2dd9c871ca543023cba
`,
		Error: errors.New("wag requires version ^1.0.0 of github.com/opentracing/opentracing-go. Please update your glide.yml and run `glide up`"),
	},
}

func TestGlideYML(t *testing.T) {
	for _, test := range glideYMLTests {
		err := ValidateGlideYML(strings.NewReader(test.YML))
		assert.Equal(t, err, test.Error,
			fmt.Sprintf("incorrect error returned from input:\n%s", test.YML))
	}
}
