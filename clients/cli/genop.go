package main

import (
	"encoding/json"
	"fmt"
	"go/types"
	"os"
	"path"
	"strings"

	. "github.com/dave/jennifer/jen"
	"github.com/getkin/kin-openapi/openapi2"
	"github.com/go-openapi/swag"
	"github.com/spf13/pflag"
)

func debugJSON(i interface{}) string {
	bs, _ := json.MarshalIndent(i, "", "  ")
	return string(bs)
}

const pkgCobra = "github.com/spf13/cobra"

func cobraCommand(s *Statement) {
	s.Qual(pkgCobra, "Command")
}

// firstArgIsContext returns whether the first argument in a function signature is context.Context.
func firstArgIsContext(sig *types.Signature) bool {
	if sig.Params().Len() < 1 {
		return false
	}
	namedType, ok := sig.Params().At(0).Type().(*types.Named)
	if !ok {
		return false
	}
	obj := namedType.Obj()
	return obj.Pkg().Name() == "context" && obj.Name() == "Context"
}

// clientMethodForOp finds the method name and signature on the Client interface
// that corresponds to the given swagger operation.
func (s swaggerExt) clientMethodForOp(op *openapi2.Operation) (string, *types.Signature, error) {
	clientMethodName := strings.Title(op.OperationID)
	var clientMethodSignature *types.Signature
	for i := 0; i < s.ClientIface.NumMethods(); i++ {
		if s.ClientIface.Method(i).Name() == clientMethodName {
			clientMethodSignature = s.ClientIface.Method(i).Type().(*types.Signature)
		}
	}
	if clientMethodSignature == nil {
		return "", nil, fmt.Errorf("could not find client method %s", clientMethodName)
	}
	return clientMethodName, clientMethodSignature, nil
}

// flagsFromParams maps swagger params to Go program flags
func flagsFromParams(params openapi2.Parameters) (*pflag.FlagSet, error) {
	fset := pflag.NewFlagSet("tmp", pflag.ContinueOnError)
	for _, param := range params {
		switch param.Type {
		case "string":
			fset.String(param.Name, "", param.Name)
		case "boolean":
			fset.Bool(param.Name, false, param.Name)
		default:
			return nil, fmt.Errorf("TODO: handle swagger param type %s", param.Type)
		}
		if param.Required {
			fset.SetAnnotation(param.Name, "required", []string{"true"})

		}
	}
	return fset, nil
}

// flagDefStatements generates the cobra flag definitions for a command.
func flagDefStatements(cmdID string, fset *pflag.FlagSet) func(s *Statement) {
	return func(s *Statement) {
		fset.VisitAll(func(flag *pflag.Flag) {
			var flagMethod string
			var flagValueLit interface{}
			switch flag.Value.Type() {
			case "string":
				flagMethod = "String"
				flagValueLit = ""
			case "bool":
				flagMethod = "Bool"
				flagValueLit = false // todo use default from swagger def
			default:
				fmt.Printf("TODO: handle flag type %s", flag.Value.Type())
				flagMethod = "TODO"
				flagValueLit = ""
			}
			s.Id(cmdID).Dot("Flags").Dot(flagMethod).Call(
				Lit(flag.Name),
				Lit(flagValueLit),
				Lit(flag.Usage),
			).Line()
			if _, ok := flag.Annotations["required"]; ok {
				s.Qual(pkgCobra, "MarkFlagRequired").Call(Id("cobraCommand").Dot("Flags").Call(), Lit(flag.Name)).Line()
			}
		})
	}
}

// cobraGetOpForBasicType returns the cobra method to retrieve a basic type, e.g. string
func cobraGetOpForBasicType(bt *types.Basic) (string, error) {
	switch bt.Name() {
	case "string":
		return "GetString", nil
	default:
		return "", fmt.Errorf("unhandled mapping from param type to cobra method: %s", paramTypeV.Name())
	}
}

// buildClientMethodInput generates the code to accumulate the input to the API client.
func buildClientMethodInput(fset *pflag.FlagSet, clientMethodName string, sig *types.Signature) (func(s *Statement), error) {
	// validate some assumptions / limitations
	if np := sig.Params().Len(); np == 0 {
		return nil, fmt.Errorf("unexpected 0 arguments: %#v", sig)
	} else if np == 1 {
		if !firstArgIsContext(sig) {
			return nil, fmt.Errorf("unexpected 1 argument that isn't context: %#v", sig)
		}
		return func(s *Statement) {}, nil // no input to construct
	} else if sig.Params().Len() != 2 {
		return nil, fmt.Errorf("unexpected > 2 arguments: %#v", sig)
	}
	// we're in the 2 argument case: context and some input type
	if !firstArgIsContext(sig) {
		return nil, fmt.Errorf("unexpected non-context first argument: %#v", sig)
	}
	param := sig.Params().At(1)
	paramType := param.Type()
	isPtr := false
	if p, ok := paramType.(*types.Pointer); ok {
		isPtr = true
		paramType = p.Elem()
	}
	inputVarName := swag.ToVarName(fmt.Sprintf("%sInput", clientMethodName))
	switch paramTypeV := paramType.(type) {
	case *types.Named:
		if !isPtr {
			return nil, fmt.Errorf("unexpected non-pointer to named type: %#v", paramTypeV)
		}
		structType, ok := paramTypeV.Underlying().(*types.Struct)
		if !ok {
			return nil, fmt.Errorf("unexpected non-struct underlying type: %#v", paramTypeV.Underlying())
		}
		obj := paramTypeV.Obj()
		fmt.Printf("client method sig pointer=%t named type: %s %s\n", isPtr, obj.Pkg().Name(), obj.Name())
		fmt.Printf("underlying: %#v\n", paramTypeV.Underlying())
		return func(s *Statement) {
			s.Id(inputVarName).Op(":=").Op("&").Qual(obj.Pkg().Path(), obj.Name()).Values(Dict{}).Line()
			for i := 0; i < structType.NumFields(); i++ {
				fieldVar := structType.Field(i)
				// TODO: figure out type of field, var it, assign it
				fmt.Printf("field %d: '%s'\n", i, fieldVar.Name())
			}
			//s.If(List(Id("v"), Err())).Op(":=").Id("cmd").
			s.Comment("TODO")
		}, nil
		//
	case *types.Basic:
		if isPtr {
			return nil, fmt.Errorf("unexpected pointer to basic type: %#v", sig)
		}
		cobraGetOp, err := cobraGetOpForBasicType(paramTypeV)
		if err != nil {
			return nil, err
		}
		return func(s *Statement) {
			s.Var().Id(inputVarName).Id(paramTypeV.Name()).Line()
			s.If(List(Id("v"), Err()).Op(":=").Id("cmd").Dot("Flags").Call().Dot(cobraGetOp).Call(Lit(param.Name())), Err().Op("!=").Nil()).Block(
				Return(Err()),
			).Else().Block(
				Id(inputVarName).Op("=").Id("v"),
			)
		}, nil

		//		fmt.Printf("client method sig %d pointer=%t basic type: %s\n", i, isPtr, paramTypeV.Name())
	default:
		return nil, fmt.Errorf("unhandled input type: %#v", paramType)
	}
	return nil, fmt.Errorf("unhandled input type: %#v", paramType)
}

// operationForPath generates code for a path specified in the spec.
func operationForPath(s *swaggerExt, pathKey string, pathValue *openapi2.PathItem) error {
	if pathValue.Get == nil {
		fmt.Println("skpping non-get", pathKey)
		return nil
	}
	fmt.Println("generating for", pathKey)
	op := pathValue.Get
	opID := op.OperationID
	cmdName := opID
	cmdFlags, err := flagsFromParams(op.Parameters)
	if err != nil {
		return err
	}
	cmdType := fmt.Sprintf("%sCmd", op.OperationID)
	clientMethodName, clientMethodSig, err := s.clientMethodForOp(op)
	if err != nil {
		return err
	}
	clientMethodInput, err := buildClientMethodInput(cmdFlags, clientMethodName, clientMethodSig)
	if err != nil {
		return err
	}

	fmt.Println(debugJSON(op))
	f := NewFilePath(path.Join(s.GoPkgGenGoCLI, "cmd"))
	//pkgJMESPath := "github.com/jmespath/go-jmespath"
	//pkgKayveeLogger := "gopkg.in/Clever/kayvee-go.v6/logger"
	//pkgYAML := "gopkg.in/yaml.v2"
	f.Comment(fmt.Sprintf("%s: %s", cmdType, op.Description))
	f.Type().Id(cmdType).Struct(
		Id("api").Func().Params(Id("cmd").Op("*").Do(cobraCommand)).Parens(List(Qual(s.GoPkgGenGoClient, "Client"), Id("error"))),
		Id("writeOutput").Func().Params(Id("cmd").Op("*").Do(cobraCommand), Id("output").Interface()).Id("error"),
	)
	f.Comment("New creates the cobra command.")
	f.Func().Parens(Id("c").Id(cmdType)).Id("New").Params().Op("*").Do(cobraCommand).Block(
		Id("cobraCmd").Op(":=").Op("&").Do(cobraCommand).Values(Dict{
			Id("Use"):   Lit(cmdName),
			Id("Short"): Lit(op.Summary),
			Id("Long"):  Lit(op.Description),
			Id("Run"):   Id("cobraRunWithContext").Call(Qual("c", "Run")),
		}),
		Do(flagDefStatements("cobraCommand", cmdFlags)),
		Return(Id("cobraCmd")),
	)
	f.Line()
	f.Comment(fmt.Sprintf("Run the %s command.", cmdName))
	f.Func().Parens(Id("c").Id(fmt.Sprintf("%s", cmdType))).Id("Run").Params(
		Id("ctx").Qual("context", "Context"),
		Id("cmd").Op("*").Do(cobraCommand),
		Id("args").Index().String(),
	).Id("error").Block(
		List(Id("api"), Err()).Op(":=").Id("c").Dot("api").Call(Id("cmd")),
		If(Err().Op("!=").Nil()).Block(Return(Err())),
		Do(clientMethodInput),
		List(Id(swag.ToVarName(fmt.Sprintf("%sOutput", clientMethodName))), Err()).Op(":=").Id("api").Dot(clientMethodName).Call(Id("ctx"), Id(swag.ToVarName(fmt.Sprintf("%sInput", clientMethodName)))),
		If(Err().Op("!=").Nil()).Block(Return(Err())),
		Return(Id("c").Dot("writeOutput").Call(Id("cmd"), Id(swag.ToVarName(fmt.Sprintf("%sOutput", clientMethodName))))),
		//Id(fmt.Sprintf("%sInput"),cmdType).Op(":=").Op("&").Qual(g.GoPkgGenGoModels, fmt.Sprintf("%s"))
	)
	f.Line()
	f.Func().Id("init").Params().Block(
		Id("rootCmd").Dot("AddCommand").Call(Id(cmdType).Values(Dict{
			Id("api"):         Id("apiFromRootFlags"),
			Id("writeOutput"): Id("writeOutputWithRootFlags"),
		}).Dot("New").Call()),
	)

	if err := os.MkdirAll(path.Join(s.GoPkgGenGoCLIPath, "cmd"), 0700); err != nil {
		return fmt.Errorf("error making directories: %s", err)
	}
	return f.Save(path.Join(s.GoPkgGenGoCLIPath, "cmd", cmdName+".go"))
}
