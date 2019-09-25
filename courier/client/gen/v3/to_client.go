package v3

import (
	"fmt"
	"io"
	"strings"

	"github.com/morlay/oas"

	"github.com/johnnyeven/libtools/codegen"
	"github.com/johnnyeven/libtools/courier/client/gen/common"
	"github.com/johnnyeven/libtools/courier/httpx"
	"github.com/johnnyeven/libtools/courier/swagger/gen"
	"github.com/johnnyeven/libtools/courier/transport_http/transform"
)

func ToClient(baseClient string, serviceName string, openAPI *oas.OpenAPI) string {
	clientSet := common.NewClientSet(baseClient, serviceName)

	for path, pathItem := range openAPI.Paths.Paths {
		for method := range pathItem.Operations.Operations {
			clientSet.AddOp(OpenAPIOperationFrom(serviceName, strings.ToUpper(string(method)), path, pathItem.Operations.Operations[method], openAPI.Components))
		}
	}

	return clientSet.String()
}

func OpenAPIOperationFrom(serviceName string, method string, path string, operation *oas.Operation, components oas.Components) *OpenAPIOperation {
	return &OpenAPIOperation{
		serviceName: serviceName,
		method:      method,
		path:        path,
		Operation:   operation,
		components:  components,
	}
}

type OpenAPIOperation struct {
	serviceName string
	method      string
	path        string
	*oas.Operation
	components oas.Components
}

var _ interface {
	common.Op
} = (*OpenAPIOperation)(nil)

func (op *OpenAPIOperation) Method() string {
	return op.method
}

func (op *OpenAPIOperation) ID() string {
	return op.Operation.OperationId
}

func (op *OpenAPIOperation) Path() string {
	return common.PathFromSwaggerPath(op.path)
}

func (op *OpenAPIOperation) HasRequest() bool {
	return len(op.Operation.Parameters) > 0 || op.RequestBody != nil
}

func (op *OpenAPIOperation) WriteReqType(w io.Writer, importer *codegen.Importer) {
	io.WriteString(w, `struct {
`)

	for _, parameter := range op.Parameters {
		schema := mayComposedFieldSchema(parameter.Schema)

		fieldName := codegen.ToUpperCamelCase(parameter.Name)
		if parameter.Extensions[gen.XField] != nil {
			fieldName = parameter.Extensions[gen.XField].(string)
		}

		field := common.NewField(fieldName)
		field.AddTag("in", string(parameter.In))
		field.AddTag("name", parameter.Name)

		field.Comment = parameter.Description

		if parameter.Extensions[gen.XTagValidate] != nil {
			field.AddTag("validate", fmt.Sprintf("%s", parameter.Extensions[gen.XTagValidate]))
		}

		if !parameter.Required {
			if schema != nil {
				d := fmt.Sprintf("%v", schema.Default)
				if schema.Default != nil && d != "" {
					field.AddTag("default", d)
				}
			}
			field.AddTag("name", parameter.Name, "omitempty")
		}

		if schema != nil {
			field.Type, _ = NewTypeGenerator(op.serviceName, importer).Type(schema)
		}

		io.WriteString(w, field.String())
	}

	if op.RequestBody != nil {
		field := common.NewField("Body")
		if jsonMedia, ok := op.RequestBody.Content[httpx.MIMEJSON]; ok && jsonMedia.Schema != nil {
			field.Type, _ = NewTypeGenerator(op.serviceName, importer).Type(jsonMedia.Schema)
			field.Comment = jsonMedia.Schema.Description
			field.AddTag("in", "body")
			field.AddTag("fmt", transform.GetContentTransformer(httpx.MIMEJSON).Key)
		}
		if formMedia, ok := op.RequestBody.Content[httpx.MIMEMultipartPOSTForm]; ok && formMedia.Schema != nil {
			field.Type, _ = NewTypeGenerator(op.serviceName, importer).Type(formMedia.Schema)
			field.Comment = formMedia.Schema.Description
			field.AddTag("in", "formData,multipart")
		}
		if formMedia, ok := op.RequestBody.Content[httpx.MIMEPOSTForm]; ok && formMedia.Schema != nil {
			field.Type, _ = NewTypeGenerator(op.serviceName, importer).Type(formMedia.Schema)
			field.Comment = formMedia.Schema.Description
			field.AddTag("in", "formData")
		}
		io.WriteString(w, field.String())
	}

	io.WriteString(w, `
}
`)
}

func (op *OpenAPIOperation) WriteRespBodyType(w io.Writer, importer *codegen.Importer) {
	respBodySchema := op.respBodySchema()
	if respBodySchema == nil {
		io.WriteString(w, `[]byte`)
		return
	}
	tpe, _ := NewTypeGenerator(op.serviceName, importer).Type(respBodySchema)
	io.WriteString(w, tpe)
}

func (op *OpenAPIOperation) respBodySchema() (schema *oas.Schema) {
	if op.Responses.Responses == nil {
		return nil
	}

	for code, resp := range op.Responses.Responses {
		if resp.Ref != "" && op.components.Responses != nil {
			if presetResponse, ok := op.components.Responses[common.RefName(resp.Ref)]; ok {
				resp = presetResponse
			}
		}

		if code >= 200 && code < 300 {
			if resp.Content[httpx.MIMEJSON] != nil {
				schema = resp.Content[httpx.MIMEJSON].Schema
				return
			}
		}
	}

	return
}
