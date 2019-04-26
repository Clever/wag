package jsclient

// PrimitiveTypes is a set of Swagger primitive types
var PrimitiveTypes = map[string]bool{
	"string":  true,
	"number":  true,
	"integer": true,
	"boolean": true,
	"array":   true,
	"object":  true,
}

// DataType represents a primitive or defined type
type DataType struct {
	Name string
	Ref  *string
}

// Param represents a function's formal parameter
type Param struct {
	Name string
	Type DataType
}
