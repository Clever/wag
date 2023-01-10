package utils

import (
	"bytes"
	"log"
	"regexp"
	"strings"
)

// ExtractModuleNameAndVersionSuffix extracts module name and suffix from a package name given an outputpath.
func ExtractModuleNameAndVersionSuffix(packageName, outputPath string) (moduleName, versionSuffix string) {
	vregex, err := regexp.Compile("/v[0-9]$|/v[0-9][0-9]")
	if err != nil {
		log.Fatalf("Error checking module name: %s", err.Error())
	}
	moduleName = strings.TrimSuffix(packageName, outputPath)
	versionSuffix = vregex.FindString(moduleName)
	moduleName = strings.TrimSuffix(moduleName, versionSuffix)
	return
}

// CamelCase converts strings to CamelCase
func CamelCase(src string, capFirst bool) string {
	camelingRegex := regexp.MustCompile("[0-9A-Za-z]+")

	// commonInitialisms, taken from
	// https://github.com/golang/lint/blob/32a87160691b3c96046c0c678fe57c5bef761456/lint.go#L702
	commonInitialisms := map[string]bool{
		"API":   true,
		"ASCII": true,
		"CPU":   true,
		"CSS":   true,
		"DNS":   true,
		"EOF":   true,
		"GUID":  true,
		"HTML":  true,
		"HTTP":  true,
		"HTTPS": true,
		"ID":    true,
		"IP":    true,
		"JSON":  true,
		"LHS":   true,
		"QPS":   true,
		"RAM":   true,
		"RHS":   true,
		"RPC":   true,
		"SLA":   true,
		"SMTP":  true,
		"SQL":   true,
		"SSH":   true,
		"TCP":   true,
		"TLS":   true,
		"TTL":   true,
		"UDP":   true,
		"UI":    true,
		"UID":   true,
		"UUID":  true,
		"URI":   true,
		"URL":   true,
		"UTF8":  true,
		"VM":    true,
		"XML":   true,
		"XSRF":  true,
		"XSS":   true,
	}

	byteSrc := []byte(src)
	chunks := camelingRegex.FindAll(byteSrc, -1)
	for idx, val := range chunks {
		if idx == 0 && !capFirst {
			chunks[idx] = val
		} else if commonInitialisms[string(bytes.ToUpper(val)[:])] {
			chunks[idx] = bytes.ToUpper(val)
		} else {
			chunks[idx] = bytes.Title(val)
		}
	}
	return string(bytes.Join(chunks, nil))
}
