package jsclient

var PrimitiveTypes = map[string]bool {
	"string": true,
	"number": true,
	"integer": true,
	"boolean": true,
	"array": true,
	"object": true,
}

type DataType struct {
	Name string
	Ref *string
}

type Param struct {
	Name string
	Type DataType
}
