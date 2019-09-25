package gin_app

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/johnnyeven/libtools/conf"
	"github.com/johnnyeven/libtools/env"
)

type GinApp struct {
	Name        string
	IP          string
	Port        int
	SwaggerPath string
	WithCORS    bool
	app         *gin.Engine
}

func (a GinApp) DockerDefaults() conf.DockerDefaults {
	return conf.DockerDefaults{
		"Port":     80,
		"WithCORS": false,
	}
}

func (a GinApp) MarshalDefaults(v interface{}) {
	if g, ok := v.(*GinApp); ok {
		if g.Name == "" {
			g.Name = os.Getenv("PROJECT_NAME")
		}

		if g.SwaggerPath == "" {
			g.SwaggerPath = "./swagger.json"
		}

		if g.Port == 0 {
			g.Port = 80
		}
	}
}

func (a *GinApp) Init() {
	if env.IsOnline() {
		gin.SetMode(gin.ReleaseMode)
	}

	a.app = gin.New()

	if a.WithCORS {
		a.app.Use(cors.New(cors.Config{
			AllowAllOrigins:  true,
			AllowMethods:     []string{"GET", "POST", "PUT", "HEAD", "DELETE", "PATCH"},
			AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "AppToken", "AccessKey"},
			AllowCredentials: false,
			MaxAge:           12 * time.Hour,
		}))
	}

	a.app.Use(gin.Recovery(), WithServiceName(a.Name), Logger())
}

type GinEngineFunc func(router *gin.Engine)

func (a *GinApp) Register(ginEngineFunc GinEngineFunc) {
	ginEngineFunc(a.app)
}

func (a *GinApp) Start() {
	a.MarshalDefaults(a)
	err := a.app.Run(a.getAddr())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Server run failed[%s]\n", err.Error())
		os.Exit(1)
	}
}

func (a GinApp) getAddr() string {
	return fmt.Sprintf("%s:%d", a.IP, a.Port)
}
