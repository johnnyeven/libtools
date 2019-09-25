package gen

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"
	"regexp"
	"strings"

	"github.com/go-openapi/spec"
	"github.com/morlay/oas"
	"github.com/sirupsen/logrus"

	"github.com/johnnyeven/libtools/codegen"
	"github.com/johnnyeven/libtools/courier/client/gen/enums"
	"github.com/johnnyeven/libtools/courier/client/gen/v2"
	"github.com/johnnyeven/libtools/courier/client/gen/v3"
	"github.com/johnnyeven/libtools/courier/status_error"
)

type ClientGenerator struct {
	File             string
	SpecURL          string
	BaseClient       string
	ServiceName      string
	swagger          *spec.Swagger
	openAPI          *oas.OpenAPI
	statusErrCodeMap status_error.StatusErrorCodeMap
}

func (g *ClientGenerator) Load(cwd string) {
	if g.SpecURL == "" && g.File == "" {
		logrus.Panicf("missing spec-url or file")
		return
	}

	if g.SpecURL != "" {
		g.loadBySpecURL()
	}

	if g.File != "" {
		g.loadByFile()
	}
}

func (g *ClientGenerator) loadByFile() {
	data, err := ioutil.ReadFile(g.File)
	if err != nil {
		panic(err)
	}

	g.swagger, g.openAPI, g.statusErrCodeMap = bytesToSwaggerOrOpenAPI(data)
}

func (g *ClientGenerator) loadBySpecURL() {
	hc := http.Client{}
	req, err := http.NewRequest("GET", g.SpecURL, nil)
	if err != nil {
		panic(err)
	}

	resp, err := hc.Do(req)
	if err != nil {
		panic(err)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	g.swagger, g.openAPI, g.statusErrCodeMap = bytesToSwaggerOrOpenAPI(bodyBytes)
	if g.ServiceName == "" {
		g.ServiceName = getIdFromUrl(g.SpecURL)
	}
}

func bytesToSwaggerOrOpenAPI(data []byte) (*spec.Swagger, *oas.OpenAPI, status_error.StatusErrorCodeMap) {
	checker := struct {
		OpenAPI string `json:"openapi"`
	}{}

	err := json.Unmarshal(data, &checker)
	if err != nil {
		panic(err)
	}

	data = bytes.Replace(data, []byte("golib/timelib"), []byte("github.com/johnnyeven/libtools/timelib"), -1)
	data = bytes.Replace(data, []byte("golib/httplib"), []byte("github.com/johnnyeven/libtools/httplib"), -1)

	statusErrCodeMap := status_error.StatusErrorCodeMap{}

	regexp.MustCompile("@httpError[^;]+;").ReplaceAllFunc(data, func(i []byte) []byte {
		v := bytes.Replace(i, []byte(`\"`), []byte(`"`), -1)
		s := status_error.ParseString(string(v))
		statusErrCodeMap[s.Code] = *s
		return i
	})

	if checker.OpenAPI == "" {
		swagger := new(spec.Swagger)

		err := json.Unmarshal(data, swagger)
		if err != nil {
			panic(err)
		}

		return swagger, nil, statusErrCodeMap
	}

	openAPI := new(oas.OpenAPI)

	err = json.Unmarshal(data, openAPI)
	if err != nil {
		panic(err)
	}

	return nil, openAPI, statusErrCodeMap
}

func getIdFromUrl(url string) string {
	paths := strings.Split(url, "/")
	return paths[len(paths)-1]
}

func (g *ClientGenerator) Pick() {
}

func (g *ClientGenerator) Output(cwd string) codegen.Outputs {
	outputs := codegen.Outputs{}
	pkgName := codegen.ToLowerSnakeCase("Client-" + g.ServiceName)

	if g.swagger != nil {
		outputs.Add(codegen.GeneratedSuffix(path.Join(pkgName, "client.go")), v2.ToClient(g.BaseClient, g.ServiceName, g.swagger))
		outputs.Add(codegen.GeneratedSuffix(path.Join(pkgName, "types.go")), v2.ToTypes(g.ServiceName, pkgName, g.swagger))
		outputs.Add(codegen.GeneratedSuffix(path.Join(pkgName, "enums.go")), enums.ToEnums(g.ServiceName, pkgName))
	}

	if g.openAPI != nil {
		outputs.Add(codegen.GeneratedSuffix(path.Join(pkgName, "client.go")), v3.ToClient(g.BaseClient, g.ServiceName, g.openAPI))
		outputs.Add(codegen.GeneratedSuffix(path.Join(pkgName, "types.go")), v3.ToTypes(g.ServiceName, pkgName, g.openAPI))
		outputs.Add(codegen.GeneratedSuffix(path.Join(pkgName, "enums.go")), enums.ToEnums(g.ServiceName, pkgName))
	}

	return outputs
}
