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

type Generator struct {
	PackageName string
	buf         bytes.Buffer
}

func (g *Generator) Printf(format string, args ...interface{}) {
	fmt.Fprintf(&g.buf, format, args...)
}

func (g *Generator) Write(p []byte) (n int, err error) {
	return g.buf.Write(p)
}

func (g *Generator) WriteFile(path string) error {
	if len(path) == 0 || path[0] == '/' {
		return fmt.Errorf("path must be relative")
	}
	fileBytes, err := format.Source(g.buf.Bytes())
	if err != nil {
		// This will error if the code isn't valid so let's print it to make it
		// easier to debug
		fmt.Printf("BAD CODE\n%s\n", string(g.buf.Bytes()))
		return err
	}
	return ioutil.WriteFile(os.Getenv("GOPATH")+"/src/"+g.PackageName+"/"+path, fileBytes, 0644)
}

// importStatements takes a list of import strings and converts them to a formatted imports
func ImportStatements(imports []string) string {
	if len(imports) == 0 {
		return ""
	}
	output := "import (\n"
	for _, importStr := range imports {
		output += fmt.Sprintf("\t\"%s\"\n", importStr)
	}
	output += ")\n\n"
	return output
}

func SortedPathItemKeys(m map[string]spec.PathItem) []string {
	sortedKeys := []string{}
	for k := range m {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)
	return sortedKeys
}

func SortedOperationsKeys(m map[string]*spec.Operation) []string {
	sortedKeys := []string{}
	for k := range m {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)
	return sortedKeys
}

func SortedStatusCodeKeys(m map[int]spec.Response) []int {
	sortedKeys := []int{}
	for k := range m {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Ints(sortedKeys)
	return sortedKeys
}

func Capitalize(input string) string {
	return strings.ToUpper(input[0:1]) + input[1:]
}

// CapitalCamelCase converts a string from snake_case into camel case. It leaves
// non-snake-case strings as is.
func SnakeToCamelCase(input string) string {
	output := ""
	parts := strings.Split(input, "_")
	for _, part := range parts {
		// Special case ID since it comes up so often and go's best practices are to
		// make it captial ID
		if part == "id" {
			output += "ID"
		} else {
			output += Capitalize(part)
		}
	}
	return output
}

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
