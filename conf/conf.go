package conf

import (
	"encoding"
	"fmt"
	"go/ast"
	"os"
	"reflect"
	"strings"

	"github.com/johnnyeven/libtools/courier/transport_http/transform"
	"github.com/johnnyeven/libtools/reflectx"
	"github.com/johnnyeven/libtools/strutil"
	"github.com/johnnyeven/libtools/validate"

	"github.com/sirupsen/logrus"
)

func UnmarshalConf(c interface{}, prefix string) EnvVars {
	rv := reflectx.Indirect(reflect.ValueOf(c))
	tpe := reflectx.IndirectType(reflect.TypeOf(c))

	if !rv.CanSet() || rv.Type().Kind() != reflect.Struct {
		panic("UnmarshalConf need an variable which can set")
	}

	ok, errMsgs := NewScanner(prefix).Unmarshal(rv, tpe)
	if !ok {
		for k, v := range errMsgs {
			logrus.Errorf("%s: %s", k, v)
		}
		logrus.Panic()
	}

	envVars, err := CollectEnvVars(rv, prefix)
	if err != nil {
		panic(err)
	}
	InitialRoot(rv)
	return envVars
}

type ICanInit interface {
	Init()
}

func InitialRoot(rv reflect.Value) {
	tpe := rv.Type()
	for i := 0; i < tpe.NumField(); i++ {
		value := rv.Field(i)
		if conf, ok := value.Interface().(ICanInit); ok {
			conf.Init()
		}
	}
}

// check and modify value
type IDefaultsMarshaller interface {
	MarshalDefaults(v interface{})
}

func NewScanner(prefix string) *Scanner {
	if prefix == "" {
		prefix = "s"
	}
	return &Scanner{
		prefix: prefix,
	}
}

type Scanner struct {
	prefix    string
	walker    transform.PathWalker
	errMsgMap transform.ErrMsgMap
}

func (vs *Scanner) Unmarshal(rv reflect.Value, tpe reflect.Type) (bool, transform.ErrMsgMap) {
	vs.marshalAndValidate(rv, tpe, "")
	return len(vs.errMsgMap) == 0, vs.errMsgMap
}

func (vs *Scanner) setErrMsg(path string, msg string) {
	if vs.errMsgMap == nil {
		vs.errMsgMap = transform.ErrMsgMap{}
	}
	vs.errMsgMap[path] = msg
}

func (vs *Scanner) getEnvKey() string {
	key := strings.ToUpper(vs.prefix)

	for _, p := range vs.walker.Paths() {
		key += strings.ToUpper(fmt.Sprintf("_%v", p))
	}

	return key
}

func (vs *Scanner) marshalAndValidate(rv reflect.Value, tpe reflect.Type, tagValidate string) {
	v := rv.Interface()
	if rv.Kind() != reflect.Ptr {
		v = rv.Addr().Interface()
	}

	if defaultsMarshaller, ok := v.(IDefaultsMarshaller); ok {
		defaultsMarshaller.MarshalDefaults(v)
	}

	if _, ok := v.(encoding.TextUnmarshaler); ok {
		errMsg := marshalEnvValueAndValidate(rv, tpe, vs.getEnvKey(), tagValidate)
		if errMsg != "" {
			vs.setErrMsg(vs.walker.String(), errMsg)
		}
		return
	}

	tpe = reflectx.IndirectType(tpe)

	switch tpe.Kind() {
	case reflect.Func:
		// skip func
	case reflect.Struct:
		if rv.Kind() == reflect.Ptr {
			if rv.IsNil() && rv.CanSet() {
				rv.Set(reflect.New(reflectx.IndirectType(tpe)))
			}
		}

		rv = reflectx.Indirect(rv)

		for i := 0; i < tpe.NumField(); i++ {
			field := tpe.Field(i)
			if !ast.IsExported(field.Name) {
				continue
			}

			if !field.Anonymous {
				vs.walker.Enter(field.Name)
			}

			tagValidate, _ := transform.GetTagValidate(&field)

			vs.marshalAndValidate(rv.Field(i), field.Type, tagValidate)

			if !field.Anonymous {
				vs.walker.Exit()
			}
		}
	default:
		errMsg := marshalEnvValueAndValidate(rv, tpe, vs.getEnvKey(), tagValidate)
		if errMsg != "" {
			vs.setErrMsg(vs.walker.String(), errMsg)
		}
	}
}

func marshalEnvValueAndValidate(
	rv reflect.Value,
	tpe reflect.Type,
	envKey string,
	tagValidate string,
) string {
	envValue, _ := os.LookupEnv(envKey)

	isPtr := rv.Kind() == reflect.Ptr

	if isPtr && rv.IsNil() {
		// initial ptr
		if rv.CanSet() {
			rv.Set(reflect.New(reflectx.IndirectType(tpe)))
		}
	}

	rv = reflectx.Indirect(rv)

	if envValue != "" && rv.CanSet() {
		err := strutil.ConvertFromStr(envValue, rv)
		if err != nil {
			return fmt.Sprintf("%s can't set wrong default value %s", rv.Type().Name(), envValue)
		}
	}

	if tagValidate != "" {
		isValid, msg := validate.ValidateItem(tagValidate, rv.Interface(), "")
		if !isValid {
			return msg
		}
	}

	return ""
}

func getConfTagFlags(tagStr string) map[string]bool {
	flagList := strings.Split(tagStr, ",")
	flags := map[string]bool{}
	for _, f := range flagList {
		flags[f] = true
	}
	return flags
}
