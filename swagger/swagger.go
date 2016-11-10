package swagger

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/go-openapi/spec"
)

// Generator handles common code generation operations when generating a file in a Go package.
type Generator struct {
	PackageName string
	buf         bytes.Buffer
}

// Printf writes a formatted string to the buffer.
func (g *Generator) Printf(format string, args ...interface{}) {
	fmt.Fprintf(&g.buf, format, args...)
}

// Write bytes to the buffer.
func (g *Generator) Write(p []byte) (n int, err error) {
	return g.buf.Write(p)
}

// WriteFile writes the buffer to a gofmt-ed file.
// The file will be located at $GOPATH/src/{PackageName}/{path}.
func (g *Generator) WriteFile(path string) error {
	if len(path) == 0 || path[0] == '/' {
		return fmt.Errorf("path must be relative")
	}
	fileBytes, err := format.Source(g.buf.Bytes())
	if err != nil {
		// This will error if the code isn't valid so let's write it out so we can debug
		f, createErr := os.Create("badcode.txt")
		if createErr != nil {
			return createErr
		}
		if _, writeErr := f.Write(g.buf.Bytes()); writeErr != nil {
			return writeErr
		}

		return fmt.Errorf("INTERNAL ERROR: %s. The invalid code was written to badcode.txt", err)
	}
	return ioutil.WriteFile(os.Getenv("GOPATH")+"/src/"+g.PackageName+"/"+path, fileBytes, 0644)
}

// ImportStatements takes a list of import strings and converts them to a formatted imports
func ImportStatements(imports []string) string {
	if len(imports) == 0 {
		return ""
	}
	remoteImports := []string{}
	output := "import (\n"
	for _, importStr := range imports {
		if strings.Contains(importStr, ".") {
			remoteImports = append(remoteImports, importStr)
		} else {
			output += fmt.Sprintf("\t\"%s\"\n", importStr)
		}
	}
	if len(remoteImports) > 0 {
		output += "\n"
	}
	for _, importStr := range remoteImports {
		output += fmt.Sprintf("\t\"%s\"\n", importStr)
	}
	output += ")\n\n"
	return output
}

// SortedPathItemKeys sorts the keys of a map[string]spec.PathItem.
func SortedPathItemKeys(m map[string]spec.PathItem) []string {
	sortedKeys := []string{}
	for k := range m {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)
	return sortedKeys
}

// SortedOperationsKeys sorts the keys of a map[string]*spec.Operation.
func SortedOperationsKeys(m map[string]*spec.Operation) []string {
	sortedKeys := []string{}
	for k := range m {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)
	return sortedKeys
}

// SortedStatusCodeKeys sorts the keys of a map[int]spec.Response.
func SortedStatusCodeKeys(m map[int]spec.Response) []int {
	sortedKeys := []int{}
	for k := range m {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Ints(sortedKeys)
	return sortedKeys
}

// SortedResponses sorts the keys of a map[string[spec].Response
func SortedResponses(m map[string]spec.Response) []string {
	sortedKeys := []string{}
	for k := range m {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)
	return sortedKeys
}

// SortedSchemaProperties sorts the properties of a schema
func SortedSchemaProperties(schema spec.Schema) []string {
	sortedKeys := []string{}
	for k := range schema.Properties {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)
	return sortedKeys
}

// Capitalize the first character of a string.
func Capitalize(input string) string {
	return strings.ToUpper(input[0:1]) + input[1:]
}

// PathItemOperations converts a spec.PathItem to a map from HTTP operation (e.g., "GET") to spec.Operation.
func PathItemOperations(item spec.PathItem) map[string]*spec.Operation {
	ops := make(map[string]*spec.Operation)
	if item.Get != nil {
		ops["GET"] = item.Get
	}
	if item.Put != nil {
		ops["PUT"] = item.Put
	}
	if item.Post != nil {
		ops["POST"] = item.Post
	}
	if item.Delete != nil {
		ops["DELETE"] = item.Delete
	}
	if item.Options != nil {
		ops["OPTIONS"] = item.Options
	}
	if item.Head != nil {
		ops["HEAD"] = item.Head
	}
	if item.Patch != nil {
		ops["PATCH"] = item.Patch
	}
	return ops
}
