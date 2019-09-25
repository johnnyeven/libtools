package dockerizier

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/johnnyeven/libtools/conf"
)

func Dockerize(envVars conf.EnvVars, serviceName string) {
	writeToFile("./dockerfile.default.yml", toDockerFileYML(envVars, serviceName))
	writeToFile("./config/default.yml", toConfigDefaultYML(envVars))
	writeToFile("./docker-compose.default.yml", toDockerComposeYML(envVars, serviceName))
}

func writeToFile(filename string, content string) error {
	dir := filepath.Dir(filename)
	if dir != "" {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return ioutil.WriteFile(filename, []byte(content), os.ModePerm)
}
