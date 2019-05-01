package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"go/parser"
	"go/types"
	"io/ioutil"
	"log"
	"os"
	"path"

	. "github.com/dave/jennifer/jen"
	"github.com/getkin/kin-openapi/jsoninfo"
	"github.com/getkin/kin-openapi/openapi2"
	"github.com/icza/dyno"
	"golang.org/x/tools/go/loader"
	yaml "gopkg.in/yaml.v2"
)

// swaggerExt adds some fields necessary to make the CLI package.
type swaggerExt struct {
	*openapi2.Swagger
	AppName           string
	GoPkg             string
	GoPkgGenGoCLI     string
	GoPkgGenGoCLIPath string
	GoPkgGenGoClient  string
	ClientIface       *types.Interface
}

func newSwaggerExt(s *openapi2.Swagger) (*swaggerExt, error) {
	e := &swaggerExt{Swagger: s}
	if s.Info.Title == "" {
		return nil, errors.New("must specify 'title' field in info")
	}
	e.AppName = s.Info.Title
	goPkgJSONRaw, ok := s.Info.Extensions["x-go-package"].(json.RawMessage)
	if !ok {
		return nil, errors.New("must specify x-go-package in info, e.g. github.com/Clever/<repo>")
	}
	var goPkg string
	if err := json.Unmarshal(goPkgJSONRaw, &goPkg); err != nil {
		return nil, err
	}
	e.GoPkg = goPkg
	e.GoPkgGenGoCLI = path.Join(e.GoPkg, "gen-go", "cli")
	e.GoPkgGenGoCLIPath = path.Join(os.Getenv("GOPATH"), "src", e.GoPkg, "gen-go", "cli")
	e.GoPkgGenGoClient = path.Join(e.GoPkg, "gen-go", "client")
	conf := loader.Config{ParserMode: parser.ParseComments}
	conf.Import(e.GoPkgGenGoClient)
	conf.Import(path.Join(e.GoPkg, "gen-go", "models"))
	lprog, err := conf.Load()
	if err != nil {
		return nil, fmt.Errorf("go/loader error loading gen-go: %s", err)
	}
	if o := lprog.Package(e.GoPkgGenGoClient).Pkg.Scope().Lookup("Client"); o == nil {
		return nil, fmt.Errorf("%s.%s not found", e.GoPkgGenGoClient, "Client")
	} else if iface, ok := o.Type().Underlying().(*types.Interface); !ok {
		return nil, fmt.Errorf("%s.%s not an interface type", e.GoPkgGenGoClient, "Client")
	} else {
		e.ClientIface = iface
	}
	return e, nil
}

func main() {
	s, err := readSwaggerYML(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	if err := genCLI(s); err != nil {
		log.Fatal(err)
	}
}

// maingo creates main.go
func maingo(s *swaggerExt) error {
	f := NewFilePathName(path.Join(s.GoPkgGenGoCLI), "main")
	f.Func().Id("main").Params().Block(
		Qual(path.Join(s.GoPkgGenGoCLI, "cmd"), "Execute").Call(),
	)
	if err := f.Save(path.Join(s.GoPkgGenGoCLIPath, "main.go")); err != nil {
		return err
	}
	return nil
}

// commongo creates cmd/common.go
func commongo(s *swaggerExt) error {
	f := NewFilePath(path.Join(s.GoPkgGenGoCLI, "cmd"))
	pkgJMESPath := "github.com/jmespath/go-jmespath"
	pkgKayveeLogger := "gopkg.in/Clever/kayvee-go.v6/logger"
	pkgYAML := "gopkg.in/yaml.v2"
	pkgCobra := "github.com/spf13/cobra"
	pkgTerminal := "golang.org/x/crypto/ssh/terminal"
	f.ImportName(pkgJMESPath, "jmespath")
	f.ImportName(pkgYAML, "yaml")
	f.Comment("writeOutputWithRootFlags outputs data using the options specified by the root command:")
	f.Comment("--query to apply a JMESPath query to the object")
	f.Comment("--format to output either json or yml")
	f.Comment("--no-pager to prevent output from being sent to pager")
	f.Func().Id("writeOutputWithRootFlags").Params(
		Id("cmd").Op("*").Qual(pkgCobra, "Command"),
		Id("output").Interface(),
	).Id("error").Block(
		If(List(Id("j"), Err()).Op(":=").Id("cmd").Dot("Flags").Call().Dot("GetString").Call(Lit("query")), Err().Op("!=").Nil()).Block(
			Return(Err()),
		).Else().If(Id("j").Op("!=").Lit("")).Block(
			List(Id("jpath"), Err()).Op(":=").Qual(pkgJMESPath, "Compile").Call(Id("j")),
			If(Err().Op("!=").Nil()).Block(
				Return(Qual("fmt", "Errorf").Call(Lit("invalid query '%s': %s"), Id("j"), Err().Dot("Error").Call())),
			),
			List(Id("bs"), Err()).Op(":=").Qual("encoding/json", "Marshal").Call(Id("output")),
			If(Err().Op("!=").Nil()).Block(
				Return(Err()),
			),
			Var().Id("outputJSONMap").Interface(),
			If(Err().Op(":=").Qual("encoding/json", "Unmarshal").Call(Id("bs"), Id("&outputJSONMap")), Err().Op("!=").Nil()).Block(
				Return(Err()),
			),
			List(Id("newOutput"), Err()).Op(":=").Id("jpath.Search").Call(Id("outputJSONMap")),
			If(Err().Op("!=").Nil()).Block(
				Return(Qual("fmt", "Errorf").Call(Lit("error applying query: %s"), Err().Dot("Error").Call())),
			),
			Id("output").Op("=").Id("newOutput"),
		),
		Line(),
		List(Id("f"), Err()).Op(":=").Id("cmd.Flags").Call().Dot("GetString").Call(Lit("format")),
		If(Err().Op("!=").Nil()).Block(Return(Err())),
		Var().Id("outputString").String(),
		Switch(Id("f")).Block(
			Case(Lit("json")).Block(
				List(Id("bs"), Err()).Op(":=").Qual("encoding/json", "MarshalIndent").Call(Id("output"), Lit(""), Lit("  ")),
				If(Err().Op("!=").Nil()).Block(
					Return(Err()),
				),
				Id("outputString").Op("=").String().Call(Id("bs")),
			),
			Case(Lit("yml")).Block(
				List(Id("bs"), Err()).Op(":=").Qual(pkgYAML, "Marshal").Call(Id("output")),
				If(Err().Op("!=").Nil()).Block(
					Return(Err()),
				),
				Id("outputString").Op("=").String().Call(Id("bs")),
			),
			Default().Block(
				Return(Qual("fmt", "Errorf").Call(Lit("unsupported format: %s"), Id("f"))),
			),
		),
		Line(),
		Comment("if output exceeds terminal size, use pager"),
		List(Id("winSize"), Err()).Op(":=").Id("getWinsize").Call(),
		If(Err().Op("!=").Nil()).Block(Return(Err())),
		List(Id("nopager"), Err()).Op(":=").Id("cmd").Dot("Flags").Call().Dot("GetBool").Call(Lit("no-pager")),
		If(Err().Op("!=").Nil()).Block(Return(Err())),
		Id("isTerminal").Op(":=").Qual(pkgTerminal, "IsTerminal").Call(Int().Call(Qual("os", "Stdout").Dot("Fd").Call())),
		If(Id("nopager").Op("||").Qual("strings", "Count").Call(Id("outputString"), Lit("\n")).Op("<=").Int().Call(Id("winSize").Dot("Row")).Op("||").Op("!").Id("isTerminal")).Block(
			Qual("fmt", "Println").Call(Id("outputString")),
		).Else().Block(
			Id("pager").Op(":=").Lit("/usr/bin/less"),
			If(Id("userPager").Op(":=").Qual("os", "Getenv").Call(Lit("PAGER")), Id("userPager").Op("!=").Lit("")).Block(
				Id("pager").Op("=").Id("userPager"),
			),
			Id("cmd").Op(":=").Qual("os/exec", "Command").Call(Id("pager")),
			Id("cmd").Dot("Stdin").Op("=").Qual("strings", "NewReader").Call(Id("outputString")),
			Id("cmd").Dot("Stdout").Op("=").Qual("os", "Stdout"),
			Err().Op(":=").Id("cmd").Dot("Run").Call(),
			If(Err().Op("!=").Nil()).Block(Return(Err())),
		),
		Return(Nil()),
	)
	f.Comment("apiFromRootFlags constructs an API client using the options specified by the")
	f.Comment("root command.")
	f.Func().Id("apiFromRootFlags").Params(Id("cmd").Op("*").Qual(pkgCobra, "Command")).Parens(List(Qual(s.GoPkgGenGoClient, "Client"), Id("error"))).Block(
		Var().Id("c").Op("*").Qual(s.GoPkgGenGoClient, "WagClient"),
		If(List(Id("env"), Err()).Op(":=").Id("cmd").Dot("Flags").Call().Dot("GetString").Call(Lit("environment")), Err().Op("!=").Nil()).Block(
			Return(List(Nil(), Err())),
		).Else().If(Id("env").Op("!=").Lit("")).Block(
			Id("c").Op("=").Id("client").Dot("New").Call(Qual("fmt", "Sprintf").Call(Lit("https://%s--"+s.AppName+".int.clever.com"), Id("env"))),
		).Else().If(List(Id("addr"), Err()).Op(":=").Id("cmd").Dot("Flags").Call().Dot("GetString").Call(Lit("addr")), Err().Op("!=").Nil()).Block(
			Return(Nil(), Err()),
		).Else().If(Id("addr").Op("!=").Lit("")).Block(
			Id("c").Op("=").Id("client").Dot("New").Call(Id("addr")),
		).Else().Block(
			Id("c").Op("=").Id("client").Dot("New").Call(Lit("http://localhost:8080")),
		),
		Line(),
		Comment("wag client produces circuit breaker logs and logs on errors, ignore them"),
		Id("nullLogger").Op(":=").Qual(pkgKayveeLogger, "New").Call(Lit("null")),
		Id("nullLogger").Dot("SetOutput").Call(Qual("io/ioutil", "Discard")),
		Id("c").Dot("SetLogger").Call(Id("nullLogger")),
		Return(Id("c"), Nil()),
	)
	f.Comment("readBodyInputWithRootFlags reads body inputs (e.g. to POST/PUT methods) from either")
	f.Comment("stdin or the user's $EDITOR if interactive is sent")
	f.Func().Id("readBodyInputWithRootFlags").Params(Id("cmd").Op("*").Qual(pkgCobra, "Command"), Id("example").Interface(), Id("into").Interface()).Id("error").Block(
		List(Id("interactive"), Err()).Op(":=").Id("cmd").Dot("Flags").Call().Dot("GetBool").Call(Lit("interactive")),
		If(Err().Op("!=").Nil()).Block(Return(Err())),
		Id("noStdin").Op(":=").Qual(pkgTerminal, "IsTerminal").Call(Int().Call(Qual("os", "Stdin").Dot("Fd").Call())),
		If(Id("noStdin").Op("&&").Op("!").Id("interactive")).Block(
			Return(Qual("fmt", "Errorf").Call(Lit("body input must be sent via stdin or interactively using --interactive"))),
		),
		If(Id("interactive")).Block(
			Comment("read from file after letting user edit"),
			Id("editor").Op(":=").Qual("os", "Getenv").Call(Lit("EDITOR")),
			If(Id("editor").Op("==").Lit("")).Block(
				Return(Qual("errors", "New").Call(Lit("please set $EDITOR to run interactively"))),
			),
			List(Id("exampleBs"), Err()).Op(":=").Qual("encoding/json", "MarshalIndent").Call(Id("example"), Lit(""), Lit("  ")),
			If(Err().Op("!=").Nil()).Block(Return(Err())),
			List(Id("tmpfile"), Err()).Op(":=").Qual("io/ioutil", "TempFile").Call(Lit(""), Lit("podSpec")),
			If(Err().Op("!=").Nil()).Block(Return(Err())),
			Defer().Qual("os", "Remove").Call(Id("tmpfile").Dot("Name").Call()),
			If(List(Id("_"), Err()).Op(":=").Qual("io", "Copy").Call(Id("tmpfile"), Qual("bytes", "NewReader").Call(Id("exampleBs"))), Err().Op("!=").Nil()).Block(Return(Err())),
			Id("args").Op(":=").Append(Qual("strings", "Split").Call(Id("editor"), Lit(" ")), Id("tmpfile").Dot("Name").Call()),
			Id("cmd").Op(":=").Qual("os/exec", "Command").Call(Id("args").Index(Lit(0)), Id("args").Index(Lit(1), Empty()).Op("...")),
			Id("cmd").Dot("Stdout").Op("=").Qual("os", "Stdout"),
			Id("cmd").Dot("Stdin").Op("=").Qual("os", "Stdin"),
			Id("cmd").Dot("Stderr").Op("=").Qual("os", "Stderr"),
			If(Err().Op(":=").Id("cmd").Dot("Run").Call(), Err().Op("!=").Nil()).Block(
				Return(Qual("fmt", "Errorf").Call(Lit("error running editor: %s"), Err().Dot("Error").Call())),
			),
			List(Id("bs"), Err()).Op(":=").Qual("io/ioutil", "ReadFile").Call(Id("tmpfile").Dot("Name").Call()),
			If(Err().Op("!=").Nil()).Block(
				Return(Qual("fmt", "Errorf").Call(Lit("error reading file: %s"), Err())),
			),
			If(Err().Op(":=").Qual("encoding/json", "Unmarshal").Call(Id("bs"), Id("into")), Err().Op("!=").Nil()).Block(
				Return(Qual("fmt", "Errorf").Call(Lit("error parsing JSON: %s"), Err())),
			),
			Return(Nil()),
		),
		Line(),
		Comment("not interactive, so read stdin"),
		List(Id("bs"), Err()).Op(":=").Qual("io/ioutil", "ReadAll").Call(Qual("os", "Stdin")),
		If(Err().Op("!=").Nil()).Block(
			Return(Qual("fmt", "Errorf").Call(Lit("error reading stdin: %s"), Err())),
		),
		If(Err().Op(":=").Qual("encoding/json", "Unmarshal").Call(Id("bs"), Id("into")), Err().Op("!=").Nil()).Block(
			Return(Qual("fmt", "Errorf").Call(Lit("error parsing JSON: %s"), Err())),
		),
		Return(Nil()),
	)
	f.Type().Id("winsize").Struct(
		Id("Row").Uint16(),
		Id("Col").Uint16(),
		Id("Xpixel").Uint16(),
		Id("Ypixel").Uint16(),
	)
	f.Line()
	f.Func().Id("getWinsize").Params().Parens(List(Op("*").Id("winsize"), Id("error"))).Block(
		Id("ws").Op(":=").Op("&").Id("winsize").Block(),
		List(Id("retCode"), Id("_"), Id("errno")).Op(":=").Qual("syscall", "Syscall").Call(
			Qual("syscall", "SYS_IOCTL"),
			Uintptr().Call(Qual("syscall", "Stdin")),
			Uintptr().Call(Qual("syscall", "TIOCGWINSZ")),
			Uintptr().Call(Qual("unsafe", "Pointer").Call(Id("ws"))),
		),
		If(Int().Call(Id("retCode")).Op("==").Lit(-1)).Block(
			Return(Nil(), Qual("fmt", "Errorf").Call(Lit("tiocbgwinsz error: %d"), Id("errno"))),
		),
		Return(Id("ws"), Nil()),
	)
	f.Line()
	f.Comment("cobraRun is what cobra expects a command to implement")
	f.Type().Id("cobraRun").Func().Params(Id("cmd").Op("*").Qual(pkgCobra, "Command"), Id("args").Index().String())
	f.Line()
	f.Comment("cobraRunCtx adds context")
	f.Type().Id("cobraRunCtx").Func().Params(Id("ctx").Qual("context", "Context"), Id("cmd").Op("*").Qual(pkgCobra, "Command"), Id("args").Index().String()).Error()
	f.Line()
	f.Comment("cobraRunWithContext converts a cobra command to one that is cancel-able with context.")
	f.Func().Id("cobraRunWithContext").Params(Id("run").Id("cobraRunCtx")).Id("cobraRun").Block(
		Return(Func().Params(Id("cmd").Op("*").Qual(pkgCobra, "Command"), Id("args").Index().String()).Block(
			Comment("create a context that lasts the length of the command runtime"),
			List(Id("ctx"), Id("cancelCtx")).Op(":=").Qual("context", "WithCancel").Call(Qual("context", "Background").Call()),
			Defer().Id("cancelCtx").Call(),
			Line(),
			Comment("cancel the context when signaled"),
			Id("sigChan").Op(":=").Make(Chan().Qual("os", "Signal"), Lit(1)),
			Qual("os/signal", "Notify").Call(Id("sigChan"), Qual("os", "Interrupt"), Qual("os", "Kill")),
			Go().Func().Params().Block(
				Op("<-").Id("sigChan"),
				Id("cancelCtx").Call(),
			).Call(),
			If(Err().Op(":=").Id("run").Call(Id("ctx"), Id("cmd"), Id("args")), Err().Op("!=").Nil()).Block(
				Qual("fmt", "Println").Call(Err()),
			),
		)),
	)

	return f.Save(path.Join(s.GoPkgGenGoCLIPath, "cmd/common.go"))
}

// genCLI generates a CLI
func genCLI(s2 *openapi2.Swagger) error {
	s, err := newSwaggerExt(s2)
	if err != nil {
		return err
	}

	fmt.Println(s.GoPkgGenGoCLIPath)
	if err := os.MkdirAll(s.GoPkgGenGoCLIPath, 0700); err != nil {
		return fmt.Errorf("error making directories: %s", err)
	}

	if err := maingo(s); err != nil {
		return err
	}

	if err := os.MkdirAll(path.Join(s.GoPkgGenGoCLIPath, "cmd"), 0700); err != nil {
		return fmt.Errorf("error making directories: %s", err)
	}

	if err := commongo(s); err != nil {
		return err
	}

	for pathKey, pathValue := range s.Paths {
		if err := operationForPath(s, pathKey, pathValue); err != nil {
			return err
		}
	}
	return nil
}

func toYAML(i interface{}) string {
	bs, err := yaml.Marshal(i)
	if err != nil {
		panic(err)
	}
	return string(bs)
}

func readSwaggerYML(fpath string) (*openapi2.Swagger, error) {
	bs, err := ioutil.ReadFile("/home/ubuntu/go/src/github.com/Clever/catapult/swagger.yml")
	if err != nil {
		return nil, err
	}
	// kin-openapi requires using jsoninfo decoder
	// need to convert to json to use this
	var i interface{}
	if err := yaml.Unmarshal(bs, &i); err != nil {
		return nil, err
	}
	i = dyno.ConvertMapI2MapS(i) // encoding/json doesn't support map[interface{}]interface{}
	jsonbs, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}

	decoder, err := jsoninfo.NewObjectDecoder(jsonbs)
	if err != nil {
		return nil, err
	}

	var swagger openapi2.Swagger
	return &swagger, decoder.DecodeStructFieldsAndExtensions(&swagger)
}
