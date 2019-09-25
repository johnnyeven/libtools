package gen

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/types"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"runtime/debug"
	"strings"

	"github.com/johnnyeven/libtools/courier/transport_http/transform"

	"github.com/morlay/oas"
	"github.com/sirupsen/logrus"
	"golang.org/x/tools/go/loader"

	"github.com/johnnyeven/libtools/codegen/loaderx"
	"github.com/johnnyeven/libtools/courier/httpx"
	"github.com/johnnyeven/libtools/courier/status_error"
	"github.com/johnnyeven/libtools/courier/transport_http"
)

func FullNameOfType(tpe reflect.Type) string {
	return fmt.Sprintf("%s.%s", tpe.PkgPath(), tpe.Name())
}

var TypeWebSocketListeners = FullNameOfType(reflect.TypeOf(transport_http.Listeners{}))
var TypeWebSocketClient = FullNameOfType(reflect.TypeOf(transport_http.WSClient{}))

func ConcatToOperation(method string, operators ...Operator) *oas.Operation {
	operation := &oas.Operation{}
	length := len(operators)
	for idx, operator := range operators {
		operator.BindOperation(method, operation, idx == length-1)
	}
	return operation
}

func NewOperatorScanner(program *loader.Program) *OperatorScanner {
	return &OperatorScanner{
		DefinitionScanner:  NewDefinitionScanner(program),
		StatusErrorScanner: NewStatusErrorScanner(program),
		program:            program,
	}
}

type OperatorScanner struct {
	*DefinitionScanner
	*StatusErrorScanner
	program   *loader.Program
	operators map[*types.TypeName]Operator
}

func (scanner *OperatorScanner) Operator(typeName *types.TypeName) *Operator {
	if typeName == nil {
		return nil
	}

	if operator, ok := scanner.operators[typeName]; ok {
		return &operator
	}

	defer func() {
		if e := recover(); e != nil {
			logrus.Errorf("scan Operator `%v` failed, panic: %s; calltrace: %s", typeName, fmt.Sprint(e), string(debug.Stack()))
		}
	}()

	if typeStruct, ok := typeName.Type().Underlying().(*types.Struct); ok {
		operator := Operator{
			ID:  typeName.Name(),
			Tag: getTagNameByPkgPath(typeName.Pkg().Path()),
		}

		scanner.bindParameterOrRequestBody(&operator, typeStruct)
		scanner.bindReturns(&operator, typeName)

		if scanner.operators == nil {
			scanner.operators = map[*types.TypeName]Operator{}
		}

		operator.Summary = docOfTypeName(typeName, scanner.program)

		scanner.operators[typeName] = operator

		return &operator
	}

	return nil
}

func getTagNameByPkgPath(pkgPath string) string {
	cwd, _ := os.Getwd()
	p, _ := build.Default.Import(pkgPath, "", build.FindOnly)
	tag, _ := filepath.Rel(cwd, p.Dir)
	i := strings.Index(tag, "routes/")
	if i >= 0 {
		tag = string([]byte(tag)[i:])
	}
	return strings.Replace(tag, "routes/", "", 1)
}

func (scanner *OperatorScanner) bindWebSocketMessages(op *Operator, schema *oas.Schema, typeVar *types.Var) {
	if strings.Contains(typeVar.Type().String(), TypeWebSocketClient) {
		for pkg, pkgInfo := range scanner.program.AllPackages {
			if pkg == typeVar.Pkg() {
				for selectExpr := range pkgInfo.Selections {
					if ident, ok := selectExpr.X.(*ast.Ident); ok {
						if pkgInfo.ObjectOf(ident) == typeVar && "Send" == selectExpr.Sel.Name {
							file := loaderx.FileOf(selectExpr, pkgInfo.Files...)
							ast.Inspect(file, func(node ast.Node) bool {
								switch node.(type) {
								case *ast.CallExpr:
									callExpr := node.(*ast.CallExpr)
									if callExpr.Fun == selectExpr {
										tpe := pkgInfo.TypeOf(callExpr.Args[0])
										subSchema := scanner.getSchemaByType(tpe.(*types.Named))
										op.AddWebSocketMessage(schema, subSchema)
										return false
									}
								}
								return true
							})
						}
					}
				}
			}
		}
	}
}

func (scanner *OperatorScanner) bindWebSocketListeners(op *Operator, typeFunc *types.Func) {
	scope := typeFunc.Scope()
	for _, name := range scope.Names() {
		n := scope.Lookup(name)
		if strings.Contains(n.Type().String(), TypeWebSocketListeners) {
			for pkg, pkgInfo := range scanner.program.AllPackages {
				if pkg == n.Pkg() {
					for selectExpr := range pkgInfo.Selections {
						if ident, ok := selectExpr.X.(*ast.Ident); ok {
							if pkgInfo.ObjectOf(ident) == n && "On" == selectExpr.Sel.Name {
								file := loaderx.FileOf(selectExpr, pkgInfo.Files...)
								ast.Inspect(file, func(node ast.Node) bool {
									switch node.(type) {
									case *ast.CallExpr:
										callExpr := node.(*ast.CallExpr)
										if callExpr.Fun == selectExpr {
											tpe := pkgInfo.TypeOf(callExpr.Args[0])
											schema := scanner.getSchemaByType(tpe.(*types.Named))
											op.AddWebSocketMessage(schema)

											params := pkgInfo.TypeOf(callExpr.Args[1]).(*types.Signature).Params()

											for i := 0; i < params.Len(); i++ {
												scanner.bindWebSocketMessages(op, schema, params.At(i))
											}
											return false
										}
									}
									return true
								})
							}
						}
					}
				}
			}
		}
	}
}

func (scanner *OperatorScanner) bindReturns(op *Operator, typeName *types.TypeName) {
	typeFunc := loaderx.MethodOf(typeName.Type().(*types.Named), "Output")

	if typeFunc != nil {
		metaData := ParseSuccessMetadata(docOfTypeName(typeFunc, scanner.program))

		loaderx.ForEachFuncResult(scanner.program, typeFunc, func(resultTypeAndValues ...types.TypeAndValue) {
			successType := resultTypeAndValues[0].Type

			if strings.Contains(successType.String(), TypeWebSocketListeners) {
				scanner.bindWebSocketListeners(op, typeFunc)
				return
			}

			if successType.String() != types.Typ[types.UntypedNil].String() {
				if op.SuccessType != nil && op.SuccessType.String() != successType.String() {
					logrus.Warnf(fmt.Sprintf("%s success result must be same struct, but got %v, already set %v", op.ID, successType, op.SuccessType))
				}
				op.SuccessType = successType
				op.SuccessStatus, op.SuccessResponse = scanner.getResponse(successType, metaData.Get("content-type"))
			}

			op.StatusErrors = scanner.StatusErrorScanner.StatusErrorsInFunc(typeFunc)
			op.StatusErrorSchema = scanner.DefinitionScanner.getSchemaByTypeString(statusErrorTypeString)
		})
	}
}

func (scanner *OperatorScanner) getResponse(tpe types.Type, contentType string) (status int, response *oas.Response) {
	response = &oas.Response{}

	if tpe.String() == "error" {
		status = http.StatusNoContent
		return
	}

	if contentType == "" {
		contentType = httpx.MIMEJSON
	}

	if pointer, ok := tpe.(*types.Pointer); ok {
		tpe = pointer.Elem()
	}

	if named, ok := tpe.(*types.Named); ok {
		{
			typeFunc := loaderx.MethodOf(named, "ContentType")
			if typeFunc != nil {
				loaderx.ForEachFuncResult(scanner.program, typeFunc, func(resultTypeAndValues ...types.TypeAndValue) {
					if resultTypeAndValues[0].IsValue() {
						contentType = getConstVal(resultTypeAndValues[0].Value).(string)
					}
				})
			}
		}

		{
			typeFunc := loaderx.MethodOf(named, "Status")
			if typeFunc != nil {
				loaderx.ForEachFuncResult(scanner.program, typeFunc, func(resultTypeAndValues ...types.TypeAndValue) {
					if resultTypeAndValues[0].IsValue() {
						status = int(getConstVal(resultTypeAndValues[0].Value).(int64))
					}
				})
			}
		}
	}

	response.AddContent(contentType, oas.NewMediaTypeWithSchema(scanner.DefinitionScanner.getSchemaByType(tpe)))

	return
}

func (scanner *OperatorScanner) bindParameterOrRequestBody(op *Operator, typeStruct *types.Struct) {
	for i := 0; i < typeStruct.NumFields(); i++ {
		var field = typeStruct.Field(i)

		if !field.Exported() {
			continue
		}

		var fieldType = field.Type()
		var fieldName = field.Name()
		var structFieldTags = reflect.StructTag(typeStruct.Tag(i))

		location, locationFlags := getTagNameAndFlags(structFieldTags.Get("in"))

		if location == "" {
			if fieldName == "Body" {
				location = "body"
			}
		}

		if location == "context" {
			continue
		}

		if field.Anonymous() {
			if typeStruct, ok := fieldType.Underlying().(*types.Struct); ok {
				scanner.bindParameterOrRequestBody(op, typeStruct)
			}
			continue
		}

		if location == "" {
			panic(fmt.Errorf("missing tag `in` for %s of %s", fieldName, op.ID))
		}

		name, flags := getTagNameAndFlags(structFieldTags.Get("name"))
		if name == "" {
			name, flags = getTagNameAndFlags(structFieldTags.Get("json"))
		}

		var param *oas.Parameter

		if location == "body" || location == "formData" {
			op.SetRequestBody(scanner.getRequestBody(fieldType, location, locationFlags["multipart"]))
			continue
		}

		if name == "" {
			panic(fmt.Errorf("missing tag `name` or `json` for parameter %s of %s", fieldName, op.ID))
		}

		param = scanner.getNonBodyParameter(name, flags, location, structFieldTags, fieldType)

		if param.Schema != nil && flags != nil && flags["string"] {
			param.Schema.Type = oas.TypeString
		}

		if styleValue, hasStyle := structFieldTags.Lookup("style"); hasStyle {
			param.AddExtension(XTagStyle, styleValue)
		}

		if fmtValue, hasFmt := structFieldTags.Lookup("fmt"); hasFmt {
			param.AddExtension(XTagFmt, fmtValue)
		}

		param = param.WithDesc(docOfTypeName(field, scanner.program))
		param.AddExtension(XField, field.Name())
		op.AddNonBodyParameter(param)
	}
}

func (scanner *OperatorScanner) getRequestBody(t types.Type, location string, isMultipart bool) *oas.RequestBody {
	reqBody := oas.NewRequestBody("", true)
	schema := scanner.DefinitionScanner.getSchemaByType(t)

	contentType := httpx.MIMEJSON

	if location == "formData" {
		if isMultipart {
			contentType = httpx.MIMEMultipartPOSTForm
		} else {
			contentType = httpx.MIMEPOSTForm
		}
	}

	reqBody.Required = true
	reqBody.AddContent(contentType, oas.NewMediaTypeWithSchema(schema))
	return reqBody
}

func (scanner *OperatorScanner) getNonBodyParameter(name string, nameFlags transform.TagFlags, location string, tags reflect.StructTag, t types.Type) *oas.Parameter {
	schema := scanner.DefinitionScanner.getSchemaByType(t)

	defaultValue, hasDefault := tags.Lookup("default")
	if hasDefault {
		schema.Default = defaultValue
	}

	required := true
	if hasOmitempty, ok := nameFlags["omitempty"]; ok {
		required = !hasOmitempty
	} else {
		// todo don't use non-default as required
		required = !hasDefault
	}

	validate, hasValidate := tags.Lookup("validate")
	if hasValidate {
		BindValidateFromValidateTagString(schema, validate)
	}

	if schema != nil && schema.Ref != "" {
		schema = oas.AllOf(
			schema,
			&oas.Schema{
				SchemaObject:   schema.SchemaObject,
				SpecExtensions: schema.SpecExtensions,
			},
		)
	}

	switch location {
	case "query":
		return oas.QueryParameter(name, schema, required)
	case "cookie":
		return oas.CookieParameter(name, schema, required)
	case "header":
		return oas.HeaderParameter(name, schema, required)
	case "path":
		return oas.PathParameter(name, schema)
	}
	return nil
}

type Operator struct {
	ID                string
	NonBodyParameters map[string]*oas.Parameter
	RequestBody       *oas.RequestBody

	StatusErrors      status_error.StatusErrorCodeMap
	StatusErrorSchema *oas.Schema

	Tag               string
	Summary           string
	SuccessType       types.Type
	SuccessStatus     int
	SuccessResponse   *oas.Response
	WebSocketMessages map[*oas.Schema][]*oas.Schema
}

func (operator *Operator) AddWebSocketMessage(schema *oas.Schema, returns ...*oas.Schema) {
	if operator.WebSocketMessages == nil {
		operator.WebSocketMessages = map[*oas.Schema][]*oas.Schema{}
	}
	operator.WebSocketMessages[schema] = append(operator.WebSocketMessages[schema], returns...)
}

func (operator *Operator) AddNonBodyParameter(parameter *oas.Parameter) {
	if operator.NonBodyParameters == nil {
		operator.NonBodyParameters = map[string]*oas.Parameter{}
	}
	operator.NonBodyParameters[parameter.Name] = parameter
}

func (operator *Operator) SetRequestBody(requestBody *oas.RequestBody) {
	operator.RequestBody = requestBody
}

func (operator *Operator) BindOperation(method string, operation *oas.Operation, last bool) {
	if operator.WebSocketMessages != nil {
		schema := oas.ObjectOf(nil)

		for msgSchema, list := range operator.WebSocketMessages {
			s := oas.ObjectOf(nil)

			s.SetProperty(typeOfSchema(msgSchema), msgSchema, false)

			if list != nil {
				sub := oas.ObjectOf(nil)
				for _, item := range list {
					sub.SetProperty(typeOfSchema(item), item, false)
				}
				schema.SetProperty("out", sub, false)
			}
			schema.SetProperty("in", s, false)
		}

		requestBody := oas.NewRequestBody("WebSocket", true)
		requestBody.AddContent(httpx.MIMEJSON, oas.NewMediaTypeWithSchema(schema))

		operation.SetRequestBody(requestBody)
		return
	}

	parameterNames := map[string]bool{}
	for _, parameter := range operation.Parameters {
		parameterNames[parameter.Name] = true
	}

	for _, parameter := range operator.NonBodyParameters {
		if !parameterNames[parameter.Name] {
			operation.Parameters = append(operation.Parameters, parameter)
		}
	}

	if operator.RequestBody != nil {
		operation.SetRequestBody(operator.RequestBody)
	}

	for code, statusError := range operator.StatusErrors {
		resp := (*oas.Response)(nil)
		if operation.Responses.Responses != nil {
			resp = operation.Responses.Responses[statusError.Status()]
		}
		statusErrors := status_error.StatusErrorCodeMap{}
		if resp != nil {
			statusErrors = pickStatusErrorsFromDoc(resp.Description)
		}
		statusErrors[code] = statusError
		resp = oas.NewResponse(statusErrors.String())
		resp.AddContent(httpx.MIMEJSON, oas.NewMediaTypeWithSchema(operator.StatusErrorSchema))
		operation.AddResponse(statusError.Status(), resp)
	}

	if last {
		operation.OperationId = operator.ID
		docs := strings.Split(operator.Summary, "\n")
		if operator.Tag != "" {
			operation.Tags = []string{operator.Tag}
		}
		operation.Summary = docs[0]
		if len(docs) > 1 {
			operation.Description = strings.Join(docs[1:], "\n")
		}
		if operator.SuccessType == nil {
			operation.Responses.AddResponse(http.StatusNoContent, &oas.Response{})
		} else {
			status := operator.SuccessStatus
			if status == 0 {
				status = http.StatusOK
				if method == http.MethodPost {
					status = http.StatusCreated
				}
			}
			if status >= http.StatusMultipleChoices && status < http.StatusBadRequest {
				operator.SuccessResponse = oas.NewResponse(operator.SuccessResponse.Description)
			}
			operation.Responses.AddResponse(status, operator.SuccessResponse)
		}
	}
}

func typeOfSchema(schema *oas.Schema) string {
	l := strings.Split(schema.Ref, "/")
	return l[len(l)-1]
}
