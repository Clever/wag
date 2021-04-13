package hardcoded

//go:generate go-bindata -nometadata -pkg hardcoded -o hardcoded.go ../_hardcoded/
//go:generate gofmt -w hardcoded.go
