package config_agent

import (
	"encoding/json"
	"github.com/johnnyeven/libtools/clients/client_configurations"
	"github.com/johnnyeven/libtools/conf"
	"github.com/johnnyeven/libtools/courier/client"
	"github.com/sirupsen/logrus"
	"reflect"
	"time"
)

const (
	DefaultHost           = "service-configurations.profzone.service.profzone.net"
	DefaultMode           = "http"
	DefaultPort           = 80
	DefaultRequestTimeout = 5
	DefaultPullInterval   = 60
	DefaultStackName      = "profzone"
	DefaultServiceName    = "service-configurations"
)

type Agent struct {
	Host               string `conf:"env"`
	Mode               string `conf:"env"`
	Port               int    `conf:"env"`
	Timeout            int64  `conf:"env"`
	PullConfigInterval int64  `conf:"env"`
	StackID            uint64 `conf:"env"`
	client             *client_configurations.ClientConfigurations
	config             interface{}
	rawConfig          []RawConfig
}

func (a *Agent) MarshalDefaults(v interface{}) {
	if _, ok := v.(*Agent); ok {
		if a.Host == "" {
			a.Host = DefaultHost
		}
		if a.Mode == "" {
			a.Mode = DefaultMode
		}
		if a.Port == 0 {
			a.Port = DefaultPort
		}
		if a.Timeout == 0 {
			a.Timeout = DefaultRequestTimeout
		}
		if a.PullConfigInterval == 0 {
			a.PullConfigInterval = DefaultPullInterval
		}
		if a.client == nil {
			c := &client_configurations.ClientConfigurations{
				Client: client.Client{
					Host:    a.Host,
					Mode:    a.Mode,
					Port:    a.Port,
					Timeout: time.Duration(a.Timeout) * time.Second,
				},
			}
			a.client = c
		}
		a.client.MarshalDefaults(a.client)
	}
}

func (a Agent) DockerDefaults() conf.DockerDefaults {
	return conf.DockerDefaults{
		"Host":               conf.RancherInternal(DefaultStackName, DefaultServiceName),
		"Mode":               DefaultMode,
		"Port":               DefaultPort,
		"Timeout":            DefaultRequestTimeout,
		"PullConfigInterval": DefaultPullInterval,
		"StackID":            0,
	}
}

func (a *Agent) Init() {
	a.client.Init()
}

func (a *Agent) BindConf(conf interface{}) {
	t := reflect.TypeOf(conf)
	if t.Kind() != reflect.Ptr {
		panic("the conf to be bind is not pointer.")
	}
	a.config = conf
}

func (a *Agent) Start() {
	if a.config == nil {
		panic("conf is not bind, please use BindConf to bind a configuration entry first.")
	}

	a.getFistRunConfig()
}

func (a *Agent) getFistRunConfig() {
	var result []byte
	var err error
	request := client_configurations.GetConfigurationsRequest{
		StackID: a.StackID,
		Size:    -1,
	}
	resp, err := a.client.GetConfigurations(request)
	if err != nil {
		result, err = a.loadConfigFromFile()
	} else {
		result, err = json.Marshal(resp.Body.Data)
		if err != nil {
			result, err = a.loadConfigFromFile()
		}
	}

	if err != nil {
		logrus.Panicf("load configuration failed, neither remote or local. err: %v", err)
	}

	err = json.Unmarshal(result, &a.rawConfig)
	if err != nil {
		logrus.Panicf("unmarshal raw configuration err: %v", err)
	}

	configMap := make(map[string]string)
	for _, config := range a.rawConfig {
		configMap[config.Key] = config.Value
	}

	jsonConfig, err := json.Marshal(configMap)
	if err != nil {
		logrus.Panicf("marshal raw configuration err: %v", err)
	}

	err = json.Unmarshal(jsonConfig, &a.config)
	if err != nil {
		logrus.Panicf("unmarshal configuration err: %v", err)
	}
}

func (a *Agent) loadConfigFromFile() ([]byte, error) {
	return nil, nil
}
