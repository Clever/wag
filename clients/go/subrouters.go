package goclient

import (
	"strings"

	"github.com/Clever/wag/v9/swagger"
	"github.com/Clever/wag/v9/templates"
	"github.com/go-openapi/spec"
	"github.com/iancoleman/strcase"
)

func subrouterOperationCode(
	s *spec.Swagger,
	op *spec.Operation,
	subrouter swagger.Subrouter,
) (string, error) {
	_, param := swagger.OperationInput(op)
	templateArgs := struct {
		ClientOperation string
		InputParamName  string
		OperationID     string
		SubrouterClient string
	}{
		ClientOperation: strings.ReplaceAll(
			swagger.ClientInterface(s, op),
			"models",
			subrouter.Key+"models",
		),
		InputParamName:  param,
		OperationID:     op.ID,
		SubrouterClient: strcase.ToLowerCamel(subrouter.Key) + "Client",
	}

	templateStr := `func (c *WagClient) {{.ClientOperation}} {
	return c.{{.SubrouterClient}}.{{pascalcase .OperationID}}(ctx, {{.InputParamName}})
}`

	return templates.WriteTemplate(templateStr, templateArgs)
}
