package mysql

import (
	"database/sql/driver"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/sirupsen/logrus"

	"golib/gorm"

	"github.com/johnnyeven/libtools/conf"
	"github.com/johnnyeven/libtools/conf/presets"
)

type MySQL struct {
	Name            string
	Host            string `conf:"upstream" validate:"@hostname"`
	Port            int
	User            string           `conf:"env" validate:"@string[1,)"`
	Password        presets.Password `conf:"env" validate:"@string[1,)"`
	Extra           string
	PoolSize        int
	ConnMaxLifetime time.Duration
	presets.Retry
	db *gorm.DB
}

func (m MySQL) DockerDefaults() conf.DockerDefaults {
	return conf.DockerDefaults{
		"Host": conf.RancherInternal("tool-dbs", m.Name),
		"Port": 3306,
	}
}

func (m MySQL) MarshalDefaults(v interface{}) {
	if mysql, ok := v.(*MySQL); ok {
		mysql.Retry.MarshalDefaults(&mysql.Retry)

		if mysql.Port == 0 {
			mysql.Port = 3306
		}

		if mysql.PoolSize == 0 {
			mysql.PoolSize = 10
		}

		if mysql.ConnMaxLifetime == 0 {
			mysql.ConnMaxLifetime = 4 * time.Hour
		}

		if mysql.Extra == "" {
			values := url.Values{}
			values.Set("charset", "utf8")
			values.Set("parseTime", "true")
			values.Set("interpolateParams", "true")
			values.Set("autocommit", "true")
			values.Set("loc", "Local")
			mysql.Extra = values.Encode()
		}
	}
}

func (m MySQL) GetConnect() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/?%s", m.User, m.Password, m.Host, m.Port, m.Extra)
}

func (m *MySQL) Connect() error {
	m.MarshalDefaults(m)
	db, err := connectMysql(m.GetConnect(), m.PoolSize, m.ConnMaxLifetime)
	if err != nil {
		return err
	}
	m.db = db
	return nil
}

func (m *MySQL) Init() {
	if m.db == nil {
		m.Do(m.Connect)
		m.db.SetLogger(&logger{})
	}
}

func (m *MySQL) Get() *gorm.DB {
	return m.db
}

type DBGetter interface {
	Get() *gorm.DB
}

type logger struct {
}

var sqlRegexp = regexp.MustCompile(`(\$\d+)|\?`)

func (l *logger) Print(values ...interface{}) {
	if len(values) > 1 {
		level := values[0]
		messages := []interface{}{}
		if level == "sql" {
			// sql
			var formatedValues []interface{}
			for _, value := range values[4].([]interface{}) {
				indirectValue := reflect.Indirect(reflect.ValueOf(value))
				if indirectValue.IsValid() {
					value = indirectValue.Interface()
					if t, ok := value.(time.Time); ok {
						formatedValues = append(formatedValues, fmt.Sprintf("'%v'", t.Format(time.RFC3339)))
					} else if b, ok := value.([]byte); ok {
						formatedValues = append(formatedValues, fmt.Sprintf("'%v'", string(b)))
					} else if r, ok := value.(driver.Valuer); ok {
						if value, err := r.Value(); err == nil && value != nil {
							formatedValues = append(formatedValues, fmt.Sprintf("'%v'", value))
						} else {
							formatedValues = append(formatedValues, "NULL")
						}
					} else {
						formatedValues = append(formatedValues, fmt.Sprintf("'%v'", value))
					}
				} else {
					formatedValues = append(formatedValues, fmt.Sprintf("'%v'", value))
				}
			}
			messages = append(messages, aurora.Red(fmt.Sprintf(sqlRegexp.ReplaceAllString(values[3].(string), "%v"), formatedValues...)))
			// duration
			messages = append(messages, fmt.Sprintf(" [%fms]", aurora.Magenta(float64(values[2].(time.Duration).Nanoseconds()/1e4)/100.0)))
		} else {
			messages = append(messages, values[2:]...)
		}

		logrus.WithField("tag", "gorm").Debug(messages...)
	}
}
