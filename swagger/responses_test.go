package swagger

import (
	"strings"
	"testing"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/loads/fmts"
	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBadReference(t *testing.T) {
	s := loadTestFile(t, "testyml/badref.yml")
	err := ValidateResponses(s)
	require.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "invalid schema reference"), err.Error())
}

func TestReferenceMissingMessageField(t *testing.T) {
	s := loadTestFile(t, "testyml/missingmessage.yml")
	err := ValidateResponses(s)
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "'message' field"), err.Error())
}

func TestErrorOnMissingTypes(t *testing.T) {
	s := loadTestFile(t, "testyml/missinginternal.yml")
	err := ValidateResponses(s)
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "must specify global"), err.Error())
}

func TestOtherRequiredField(t *testing.T) {
	s := loadTestFile(t, "testyml/requiredfield.yml")
	err := ValidateResponses(s)
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "cannot have required fields"), err.Error())
}

func Test3xxError(t *testing.T) {
	s := loadTestFile(t, "testyml/3xxresponse.yml")
	err := ValidateResponses(s)
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "cannot define 3XX status codes"), err.Error())
}

func TestMultiSuccessError(t *testing.T) {
	s := loadTestFile(t, "testyml/multisuccess.yml")
	err := ValidateResponses(s)
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "one success type"), err.Error())
}

func loadTestFile(t *testing.T, filename string) spec.Swagger {
	loads.AddLoader(fmts.YAMLMatcher, fmts.YAMLDoc)
	doc, err := loads.Spec(filename)
	require.NoError(t, err)
	return *doc.Spec()
}

func TestAddingDefaultTypes(t *testing.T) {
	s := loadTestFile(t, "testyml/defaults.yml")
	assert.NoError(t, ValidateResponses(s))

	responses := s.Paths.Paths["/path"].Get.Responses.StatusCodeResponses
	require.Equal(t, 3, len(responses))
}

func TestOverrideDefaults(t *testing.T) {
	s := loadTestFile(t, "testyml/override.yml")
	assert.NoError(t, ValidateResponses(s))
}
