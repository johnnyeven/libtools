package conf

import (
	"go/ast"
	"reflect"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/johnnyeven/libtools/strutil"
)

var TagConf = "conf"

type TagConfOption struct {
	CanConfig     bool `opt:"env"`
	IsUpstream    bool `opt:"upstream"`
	FallbackValue interface{}
}

func GetTagConfOption(opt string) (tagConfOption TagConfOption) {
	if opt == "" {
		return
	}

	optionList := strings.Split(opt, ",")
	options := make(map[string]bool)
	for _, opt := range optionList {
		options[opt] = true
	}

	rv := reflect.Indirect(reflect.ValueOf(&tagConfOption))
	tpe := rv.Type()

	for i := 0; i < tpe.NumField(); i++ {
		field := tpe.Field(i)
		if options[field.Tag.Get("opt")] {
			rv.Field(i).SetBool(true)
		}
	}

	return
}

type IHasDockerDefaults interface {
	DockerDefaults() DockerDefaults
}

type DockerDefaults map[string]interface{}

func (dockerDefaults DockerDefaults) Merge(nextDockerDefaults DockerDefaults) DockerDefaults {
	finalDockerDefaults := DockerDefaults{}
	for key, value := range dockerDefaults {
		finalDockerDefaults[key] = value
	}
	for key, value := range nextDockerDefaults {
		finalDockerDefaults[key] = value
	}
	return finalDockerDefaults
}

func collectEnvVars(rv reflect.Value, envVarKey string, tagConfOption TagConfOption, dockerDefaults DockerDefaults) (envVars EnvVars, err error) {
	envVars = EnvVars{}
	if rv.Kind() == reflect.Ptr && rv.IsNil() {
		return
	}
	if hasDockerDefaults, ok := rv.Interface().(IHasDockerDefaults); ok {
		// parent should overwrite child as anonymous
		dockerDefaults = hasDockerDefaults.DockerDefaults().Merge(dockerDefaults)
	}
	rv = reflect.Indirect(rv)
	switch rv.Kind() {
	case reflect.Func:
		// skip func
	case reflect.Struct:
		walkStructField(rv, false, func(fieldValue reflect.Value, field reflect.StructField) bool {
			nextEnvKey := resolveEnvVarKeyByField(envVarKey, field)
			tagContent := field.Tag.Get(TagConf)
			if tagContent == "-" {
				return true
			}
			tagConfOption = GetTagConfOption(tagContent)
			if v, exists := dockerDefaults[field.Name]; exists {
				tagConfOption.FallbackValue = v
			}
			subEnvVars, errForCollect := collectEnvVars(fieldValue, nextEnvKey, tagConfOption, dockerDefaults)
			envVars.Merge(subEnvVars)
			return errForCollect == nil
		})
	default:
		envVars.Set(envVarKey, EnvVar{
			TagConfOption: tagConfOption,
			Value:         rv,
		})
	}
	return
}

func walkStructField(rv reflect.Value, forSet bool, fn func(fieldValue reflect.Value, field reflect.StructField) bool) {
	reflectType := rv.Type()
	for i := 0; i < reflectType.NumField(); i++ {
		field := reflectType.Field(i)
		fieldValue := rv.Field(i)
		if !ast.IsExported(field.Name) {
			continue
		}

		if forSet && (!fieldValue.IsValid() || !fieldValue.CanSet()) {
			continue
		}

		next := fn(fieldValue, field)
		if !next {
			break
		}
	}
}

func resolveEnvVarKeyByField(pre string, field reflect.StructField) string {
	if field.Anonymous {
		return pre
	}
	if pre == "" {
		return strings.ToUpper(field.Name)
	}
	return strings.ToUpper(pre + "_" + field.Name)
}

func CollectEnvVars(rv reflect.Value, envVarKey string) (envVars EnvVars, err error) {
	return collectEnvVars(rv, envVarKey, TagConfOption{}, DockerDefaults{})
}

type ISecurityStringer interface {
	SecurityString() string
}

type EnvVars map[string]EnvVar

type EnvVar struct {
	Value reflect.Value
	TagConfOption
}

func stringValueOf(rv reflect.Value, security bool) string {
	v := rv.Interface()
	switch rv.Kind() {
	case reflect.Array, reflect.Slice:
		values := make([]string, 0)
		for i := 0; i < rv.Len(); i++ {
			values = append(values, stringValueOf(rv.Index(i), security))
		}
		return strings.Join(values, ",")
	default:
		if security {
			if securityStringer, ok := v.(ISecurityStringer); ok {
				return securityStringer.SecurityString()
			}
		}
		s, err := strutil.ConvertToStr(v)
		if err != nil {
			panic(err)
		}
		return s
	}
}

func (envVar EnvVar) GetFallbackValue(security bool) string {
	return stringValueOf(reflect.ValueOf(envVar.FallbackValue), security)
}

func (envVar EnvVar) GetValue(security bool) string {
	rv := envVar.Value
	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		values := make([]string, 0)
		for i := 0; i < rv.Len(); i++ {
			values = append(values, stringValueOf(rv.Index(i), security))
		}
		return strings.Join(values, ",")
	default:
		return stringValueOf(rv, security)
	}
}

func (envVars EnvVars) Set(key string, envVar EnvVar) {
	envVars[key] = envVar
}

func (envVars EnvVars) Merge(nextEnvVars EnvVars) {
	for key, envVar := range nextEnvVars {
		envVars.Set(key, envVar)
	}
}

func (envVars EnvVars) Print() {
	keysWarning := make([]string, 0)
	keysNormal := make([]string, 0)

	for key, envVar := range envVars {
		if envVar.CanConfig {
			keysWarning = append(keysWarning, key)
		} else {
			keysNormal = append(keysNormal, key)
		}
	}

	sort.Strings(keysWarning)
	sort.Strings(keysNormal)

	fields := logrus.Fields{}

	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{}

	getLogger := func(key string) *logrus.Entry {
		envVar := envVars[key]
		fields[key] = envVar.GetValue(true)
		l := logger.WithField(key, fields[key])
		if envVar.FallbackValue != nil {
			l = l.WithField("fallback", envVars[key].GetFallbackValue(true))
		}
		return l
	}

	for _, key := range keysWarning {
		getLogger(key).Warning()
	}
	for _, key := range keysNormal {
		getLogger(key).Info()
	}

	logger.Formatter = &logrus.JSONFormatter{}
	logger.WithField("config", fields).Warning()
}
