package service

import (
	"fmt"
	"os"
	"reflect"

	"github.com/spf13/cobra"

	"github.com/johnnyeven/libtools/conf"
	"github.com/johnnyeven/libtools/httplib"
	"github.com/johnnyeven/libtools/service/dockerizier"
)

// @deprecated
func New(serviceName string) *Service {
	serve := Service{Name: serviceName}
	return &serve
}

type Service struct {
	Name               string
	AutoMigrate        bool
	envConfigPrefix    string
	outputDockerConfig bool
	help               bool
	cfg                interface{}
}

func (s *Service) SetEnvConfigPrefix(pre string) {
	s.envConfigPrefix = pre
}

func (s *Service) ConfP(c interface{}) {
	tpe := reflect.TypeOf(c)
	if tpe.Kind() != reflect.Ptr {
		panic(fmt.Errorf("ConfP pass ptr for setting value"))
	}
	s.cfg = c
	s.Execute()
}

func (s *Service) Execute() {
	if projectFeature, exists := os.LookupEnv("PROJECT_FEATURE"); exists && projectFeature != "" {
		s.Name = s.Name + "--" + projectFeature
	}

	command := &cobra.Command{
		Use:   s.Name,
		Short: fmt.Sprintf("%s", s.Name),
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	command.PersistentFlags().
		StringVarP(&s.envConfigPrefix, "env-config-prefix", "e", "S", "prefix for env var config")
	command.PersistentFlags().
		BoolVarP(&s.outputDockerConfig, "output-docker-config", "c", true, "output configuration of docker")
	command.PersistentFlags().
		BoolVarP(&s.AutoMigrate, "auto-migrate", "m", os.Getenv("GOENV") == "DEV" || os.Getenv("GOENV") == "TEST", "auto migrate database if need")
	command.PersistentFlags().
		BoolVarP(&s.help, "help", "h", false, "show help")

	if err := command.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if s.help {
		os.Exit(0)
	}

	s.conf()
}

func (s *Service) conf() {
	os.Setenv(EnvVarKeyProjectName, s.Name)
	httplib.SetServiceName(s.Name)

	envVars := conf.UnmarshalConf(s.cfg, s.envConfigPrefix)
	envVars.Print()

	if s.outputDockerConfig {
		dockerizier.Dockerize(envVars, s.Name)
	}
}

var (
	EnvVarKeyProjectName = "PROJECT_NAME"
)
