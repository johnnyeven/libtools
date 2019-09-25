package log

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/johnnyeven/libtools/conf"
	"github.com/johnnyeven/libtools/log/context"
	"github.com/johnnyeven/libtools/log/hooks"
)

type Log struct {
	Name   string
	Path   string
	Level  string `conf:"env"`
	Format string
	init   bool
}

func (log Log) DockerDefaults() conf.DockerDefaults {
	return conf.DockerDefaults{
		"Format": "json",
	}
}

func (log Log) MarshalDefaults(v interface{}) {
	if l, ok := v.(*Log); ok {
		if l.Name == "" {
			l.Name = os.Getenv("PROJECT_NAME")
		}

		if l.Level == "" {
			l.Level = "DEBUG"
		}

		if l.Format == "" {
			l.Format = "text"
		}
	}
}

func (log *Log) Init() {
	if !log.init {
		log.Create()
		log.init = true
	}
}

func (log *Log) Create() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(getLogLevel(log.Level))
	if log.Format == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{
			ForceColors: true,
		})
	}
	logrus.AddHook(context.NewLogIDHook())

	logrus.AddHook(hooks.NewCallStackHook())
	logrus.AddHook(hooks.NewProjectHook(log.Name))

	logrus.SetOutput(ioutil.Discard)

	if log.Path != "" {
		logrus.AddHook(hooks.NewLogWriterHook(log.Path))
		logrus.AddHook(hooks.NewLogWriterForErrorHook(log.Path))
	} else {
		logrus.AddHook(hooks.NewLogPrinterHook())
		logrus.AddHook(hooks.NewLogPrinterForErrorHook())
	}

}

func getLogLevel(l string) logrus.Level {
	level, err := logrus.ParseLevel(strings.ToLower(l))
	if err == nil {
		return level
	}
	return logrus.InfoLevel
}
