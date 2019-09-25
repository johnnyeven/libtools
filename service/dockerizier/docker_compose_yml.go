package dockerizier

import (
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/johnnyeven/libtools/conf"
	"github.com/johnnyeven/libtools/docker"
)

func toDockerComposeYML(envVars conf.EnvVars, serviceName string) string {
	d := docker.NewDockerCompose()
	s := docker.NewService(Image)

	keys := make([]string, 0)
	for key := range envVars {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	upstreams := make([]string, 0)

	for _, key := range keys {
		envVar := envVars[key]
		if envVar.CanConfig {
			s = s.SetEnvironment(key, `${`+key+`}`)
		}
		strValue := envVar.GetValue(false)
		if strValue == "./swagger.json" {
			s = s.SetLabel("lb.g7pay.expose80", "/"+getBaseURL(serviceName))
			s = s.SetLabel("base_path", "/"+getBaseURL(serviceName))
		}

		if envVar.FallbackValue != nil {
			if envVar.IsUpstream {
				if envVar.CanConfig {
					upstreams = append(upstreams, `${`+key+`}`)
				} else {
					upstreams = append(upstreams, envVar.GetFallbackValue(false))
				}
			}
		}
	}

	if len(upstreams) > 0 {
		s = s.SetLabel("upstreams", strings.Join(upstreams, ","))
	}

	s = s.SetEnvironment("GOENV", `${GOENV}`)

	s = s.SetLabel("io.rancher.container.pull_image", "always")
	s = s.SetLabel("io.rancher.container.start_once", "true")

	s = s.SetLabel("project.description", "${PROJECT_DESCRIPTION}")
	s = s.SetLabel("project.group", "${PROJECT_GROUP}")
	s = s.SetLabel("project.name", "${PROJECT_NAME}")
	s = s.SetLabel("project.version", "${PROJECT_VERSION}")

	s = s.AddDns("169.254.169.250", "rancher.internal")

	d = d.AddService(serviceName, s)
	bytes, _ := yaml.Marshal(d)
	return string(bytes)
}

func getBaseURL(name string) string {
	pathReg := regexp.MustCompile("service-([xi]-)?(.+)")
	return pathReg.ReplaceAllStringFunc(name, func(s string) string {
		result := pathReg.FindAllStringSubmatch(s, -1)
		if len(result[0]) == 3 {
			return result[0][2]
		}
		return result[0][1]
	})
}
