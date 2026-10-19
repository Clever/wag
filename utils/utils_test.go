package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetWagVersion(t *testing.T) {
	defer SetWagVersion("")

	tests := []struct {
		name     string
		version  string
		expected string
	}{
		{name: "release version", version: "v9.11.0", expected: "v9.11.0"},
		{name: "pre-release version", version: "v9.12.0-rc.1", expected: "v9.12.0-rc.1"},
		{name: "spaces allowed", version: "v9.11.0 dev", expected: "v9.11.0 dev"},
		{name: "empty", version: "", expected: "unknown"},
		{name: "only disallowed characters", version: "*/\n", expected: "unknown"},
		{name: "surrounding whitespace trimmed", version: "  v9.11.0\n", expected: "v9.11.0"},
		{name: "newline and comment injection", version: "v9.11.0\n*/ package evil //", expected: "v9.11.0 package evil"},
		{name: "build metadata", version: "v9.11.0+abc.123", expected: "v9.11.0+abc.123"},
		{name: "build metadata with underscore", version: "v9.11.0+abc_123", expected: "v9.11.0+abc_123"},
		{name: "non-ASCII letters", version: "v9.11.0é", expected: "v9.11.0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetWagVersion(tt.version)
			assert.Equal(t, tt.expected, wagVersion)
		})
	}
}
