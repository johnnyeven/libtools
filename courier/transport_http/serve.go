package transport_http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/julienschmidt/httprouter"

	"github.com/johnnyeven/libtools/conf"
	"github.com/johnnyeven/libtools/courier"
	"os/signal"
	"github.com/johnnyeven/libtools/servicex"
)

type ServeHTTP struct {
	Name         string
	IP           string
	Port         int
	SwaggerPath  string
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
	WithCORS     bool
	router       *httprouter.Router
	serv         *http.Server
}

func (s ServeHTTP) DockerDefaults() conf.DockerDefaults {
	return conf.DockerDefaults{
		"Port":     80,
		"WithCORS": false,
	}
}

func (s ServeHTTP) MarshalDefaults(v interface{}) {
	if h, ok := v.(*ServeHTTP); ok {
		if h.Name == "" {
			h.Name = os.Getenv(servicex.EnvVarKeyProjectName)
		}

		if h.SwaggerPath == "" {
			h.SwaggerPath = "./swagger.json"
		}

		if h.Port == 0 {
			h.Port = 80
		}

		if h.ReadTimeout == 0 {
			h.ReadTimeout = 15 * time.Second
		}

		if h.WriteTimeout == 0 {
			h.WriteTimeout = 15 * time.Second
		}
	}
}

func (s *ServeHTTP) Serve(router *courier.Router) error {
	s.MarshalDefaults(s)
	s.router = s.convertRouterToHttpRouter(router)

	s.serv = &http.Server{
		Handler:      s,
		Addr:         fmt.Sprintf("%s:%d", s.IP, s.Port),
		WriteTimeout: s.WriteTimeout,
		ReadTimeout:  s.ReadTimeout,
	}

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		if err := s.Stop(); err != nil {
			fmt.Printf("HTTP server Shutdown: %v", err)
		}
	}()

	fmt.Printf("[Courier] listen on %s\n", s.serv.Addr)
	return s.serv.ListenAndServe()
}

func (s *ServeHTTP) Stop() error {
	return s.serv.Shutdown(context.Background())
}

var RxHttpRouterPath = regexp.MustCompile("/:([^/]+)")

func (s *ServeHTTP) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if s.WithCORS {
		headers := w.Header()
		setCORS(&headers, req)
	}
	s.router.ServeHTTP(w, req)
}

func (s *ServeHTTP) convertRouterToHttpRouter(router *courier.Router) *httprouter.Router {
	routes := router.Routes()

	if len(routes) == 0 {
		panic(fmt.Sprintf("need to register operation before Listion"))
	}

	r := httprouter.New()

	sort.Slice(routes, func(i, j int) bool {
		return getPath(routes[i]) < getPath(routes[j])
	})

	for _, route := range routes {
		method := getMethod(route)
		p := getPath(route)

		finalOperators, operatorTypeNames := route.EffectiveOperators()

		if len(finalOperators) == 0 {
			panic(fmt.Errorf(
				"[Courier] No available operator %v",
				route.Operators,
			))
		}

		if method == "" {
			panic(fmt.Errorf(
				"[Courier] Missing method of %s\n",
				color.CyanString(reflect.TypeOf(finalOperators[len(finalOperators)-1]).Name()),
			))
		}

		lengthOfOperatorTypes := len(operatorTypeNames)

		for i := range operatorTypeNames {
			if i < lengthOfOperatorTypes-1 {
				operatorTypeNames[i] = color.MagentaString(operatorTypeNames[i])
			} else {
				operatorTypeNames[i] = color.CyanString(operatorTypeNames[i])
			}
		}

		fmt.Printf(
			"[Courier] %s %s\n",
			colorByMethod(method)("%s %s", method[0:3], RxHttpRouterPath.ReplaceAllString(p, "/{$1}")),
			strings.Join(operatorTypeNames, " "),
		)

		r.Handle(method, p, CreateHttpHandler(s, finalOperators...))
	}

	return r
}

func getMethod(route *courier.Route) string {
	if withHttpMethod, ok := route.Operators[len(route.Operators)-1].(IMethod); ok {
		return string(withHttpMethod.Method())
	}
	return ""
}

func getPath(route *courier.Route) string {
	p := "/"
	for _, operator := range route.Operators {
		if WithHttpPath, ok := operator.(IPath); ok {
			p += WithHttpPath.Path()
		}
	}
	return httprouter.CleanPath(p)
}

func colorByMethod(method string) func(f string, args ...interface{}) string {
	switch method {
	case http.MethodGet:
		return color.BlueString
	case http.MethodPost:
		return color.GreenString
	case http.MethodPut:
		return color.YellowString
	case http.MethodDelete:
		return color.RedString
	case http.MethodHead:
		return color.WhiteString
	default:
		return color.BlackString
	}
}
