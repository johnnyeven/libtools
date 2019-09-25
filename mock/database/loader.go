package database

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

func LoadAndParse(path string, v interface{}) error {
	var content []byte
	file := path + ".json"
	_, err := os.Stat(file)
	if err == nil {
		content, err = ioutil.ReadFile(file)
		if err != nil {
			return err
		}

		return json.Unmarshal(content, v)
	}

	file = path + ".yml"
	_, err = os.Stat(file)
	if err == nil {
		content, err = ioutil.ReadFile(file)
		if err != nil {
			return err
		}

		return yaml.Unmarshal(content, v)
	}

	return fmt.Errorf("only support .json or .yml")
}
