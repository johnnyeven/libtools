package v2

import (
	"fmt"
	"io"
	"strings"

	"github.com/go-openapi/spec"

	"github.com/johnnyeven/libtools/codegen"
	"github.com/johnnyeven/libtools/courier/client/gen/common"
	"github.com/johnnyeven/libtools/courier/httpx"
	"github.com/johnnyeven/libtools/courier/swagger/gen"
)

func ToClient(baseClient string, serviceName string, swagger *spec.Swagger) string {
	clientSet := common.NewClientSet(baseClient, serviceName)

	for path, pathItem := range swagger.Paths.Paths {
		if pathItem.PathItemProps.Get != nil {
			clientSet.AddOp(SwaggerOperationFrom(serviceName, "GET", path, pathItem.PathItemProps.Get, swagger.Responses))
		}
		if pathItem.PathItemProps.Post != nil {
			clientSet.AddOp(SwaggerOperationFrom(serviceName, "POST", path, pathItem.PathItemProps.Post, swagger.Responses))
		}
		if pathItem.PathItemProps.Put != nil {
			clientSet.AddOp(SwaggerOperationFrom(serviceName, "PUT", path, pathItem.PathItemProps.Put, swagger.Responses))
		}
		if pathItem.PathItemProps.Delete != nil {
			clientSet.AddOp(SwaggerOperationFrom(serviceName, "DELETE", path, pathItem.PathItemProps.Delete, swagger.Responses))
		}
		if pathItem.PathItemProps.Head != nil {
			clientSet.AddOp(SwaggerOperationFrom(serviceName, "HEAD", path, pathItem.PathItemProps.Head, swagger.Responses))
		}
		if pathItem.PathItemProps.Patch != nil {
			clientSet.AddOp(SwaggerOperationFrom(serviceName, "PATCH", path, pathItem.PathItemProps.Patch, swagger.Responses))
		}
		if pathItem.PathItemProps.Options != nil {
			clientSet.AddOp(SwaggerOperationFrom(serviceName, "OPTIONS", path, pathItem.PathItemProps.Options, swagger.Responses))
		}
	}

	return clientSet.String()
}

func SwaggerOperationFrom(serviceName string, method string, path string, operation *spec.Operation, responses map[string]spec.Response) *SwaggerOperation {
	return &SwaggerOperation{
		serviceName: serviceName,
		method:      method,
		path:        path,
		Operation:   operation,
		responses:   responses,
	}
}

type SwaggerOperation struct {
	serviceName string
	method      string
	path        string
	*spec.Operation
	responses map[string]spec.Response
}

var _ interface {
	common.Op
} = (*SwaggerOperation)(nil)

func (op *SwaggerOperation) Method() string {
	return op.method
}

func (op *SwaggerOperation) ID() string {
	return op.Operation.ID
}

func (op *SwaggerOperation) Path() string {
	return common.PathFromSwaggerPath(op.path)
}

func (op *SwaggerOperation) HasRequest() bool {
	return len(op.Operation.Parameters) > 0
}

func (op *SwaggerOperation) WriteReqType(w io.Writer, importer *codegen.Importer) {
	io.WriteString(w, `struct {
`)

	for _, parameter := range op.Parameters {
		fieldName := codegen.ToUpperCamelCase(parameter.Name)
		if parameter.Extensions[gen.XField] != nil {
			fieldName = parameter.Extensions[gen.XField].(string)
		}

		field := common.NewField(fieldName)
		field.AddTag("in", parameter.In)

		tagName := parameter.Name
		if parameter.Extensions[gen.XTagName] != nil {
			tagName = fmt.Sprintf("%s", parameter.Extensions[gen.XTagName])
		}
		flags := make([]string, 0)
		if !parameter.Required && !strings.Contains(tagName, "omitempty") {
			flags = append(flags, "omitempty")
		}
		field.AddTag("name", tagName, flags...)
		field.Comment = parameter.Description

		if parameter.Extensions[gen.XTagValidate] != nil {
			field.AddTag("validate", fmt.Sprintf("%s", parameter.Extensions[gen.XTagValidate]))
		}

		if parameter.Default != nil {
			defaultValue := fmt.Sprintf("%v", parameter.Default)
			if defaultValue != "" {
				field.AddTag("default", defaultValue)
			}
		}

		if parameter.Schema != nil {
			field.Type = NewTypeGenerator(op.serviceName, importer).Type(parameter.Schema)
		} else {
			schema := spec.Schema{}
			schema.Typed(parameter.Type, parameter.Format)
			if parameter.Items != nil {
				itemSchema := spec.Schema{}
				itemSchema.Typed(parameter.Items.Type, parameter.Items.Format)
				itemSchema.VendorExtensible = parameter.Items.VendorExtensible
				schema.Items = &spec.SchemaOrArray{
					Schema: &itemSchema,
				}
			}
			schema.Extensions = parameter.Extensions
			field.Type = NewTypeGenerator(op.serviceName, importer).Type(&schema)
		}

		io.WriteString(w, field.String())
	}

	io.WriteString(w, `
}
`)
}

func (op *SwaggerOperation) WriteRespBodyType(w io.Writer, importer *codegen.Importer) {
	respBodySchema := op.respBodySchema()
	if respBodySchema == nil {
		io.WriteString(w, `[]byte`)
		return
	}
	io.WriteString(w, NewTypeGenerator(op.serviceName, importer).Type(respBodySchema))
}

func (op *SwaggerOperation) respBodySchema() *spec.Schema {
	if op.Responses == nil || op.Responses.StatusCodeResponses == nil {
		return nil
	}

	var hasStringProduce = false

	for _, produce := range op.Produces {
		hasStringProduce = produce == httpx.MIMEHTML
	}

	var schema *spec.Schema

	for code, r := range op.Responses.StatusCodeResponses {
		if code >= 200 && code < 300 {
			schema = &spec.Schema{}

			if hasStringProduce {
				schema.Typed("string", "")
			} else {
				schema.Typed("null", "")
			}

			if r.Ref.String() != "" && op.responses != nil {
				if op.responses[common.RefName(r.Ref.String())].Schema != nil {
					schema = op.responses[common.RefName(r.Ref.String())].Schema
				}
			}

			if r.Schema != nil {
				schema = r.Schema
			}
		}
	}

	return schema
}
