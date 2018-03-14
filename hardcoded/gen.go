package hardcoded

//go:generate $PWD/bin/go-bindata -pkg hardcoded -o hardcoded.go ../_hardcoded/
//go:generate gofmt -w hardcoded.go
