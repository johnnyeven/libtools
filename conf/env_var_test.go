package conf

import (
	"reflect"
	"testing"

	"github.com/johnnyeven/libtools/ptr"

	"github.com/stretchr/testify/assert"

	"github.com/johnnyeven/libtools/conf/presets"
)

type SubConfig struct {
	Key  string `conf:"env"`
	Bool bool
	Func func() error
}

func (s SubConfig) DockerDefaults() DockerDefaults {
	return DockerDefaults{
		"Key":  "test",
		"Bool": false,
	}
}

func TestEnvVar(t *testing.T) {
	tt := assert.New(t)

	c := struct {
		Array     []string
		Password  presets.Password `conf:"env"`
		PtrString *string          `conf:"env"`
		Host      string           `validate:"@hostname"`
		SubConfig
	}{}

	c.Password = "123456"
	c.Key = "123456"
	c.PtrString = ptr.String("123456")
	c.Array = []string{"1", "2"}
	c.Password = "host/"

	rv := reflect.Indirect(reflect.ValueOf(c))

	envVars, err := CollectEnvVars(rv, "T")
	tt.Nil(err)

	envVars.Print()
}
