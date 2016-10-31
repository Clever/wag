package jsclient

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFillOutPath(t *testing.T) {
	testSpecs := []struct {
		i string
		o string
	}{
		{"/url", "/url"},
		{"/url/{param}", `/url/" + params.param + "`},
		{"/url/{Param}", `/url/" + params.Param + "`},
		{"/url/{param}/other/{longParam}", `/url/" + params.param + "/other/" + params.longParam + "`},
		{"/url/{param_1}/other/{param_2}", `/url/" + params.param1 + "/other/" + params.param2 + "`},
	}
	for _, spec := range testSpecs {
		assert.Equal(t, spec.o, fillOutPath(spec.i))
	}
}
