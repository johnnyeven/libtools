package redis

import (
	"fmt"
	"os"
	"time"

	"github.com/johnnyeven/libtools/conf"
	"github.com/johnnyeven/libtools/conf/presets"
	"github.com/johnnyeven/libtools/env"
)

type Redis struct {
	Name           string
	Protocol       string
	Host           string `conf:"upstream" validate:"@hostname"`
	Port           int
	Password       presets.Password `conf:"env" validate:"@string(0,)"`
	ConnectTimeout time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	IdleTimeout    time.Duration
	MaxActive      int
	MaxIdle        int
	Wait           bool
	DB             string
	cache          *RedisCache
}

func (r Redis) DockerDefaults() conf.DockerDefaults {
	return conf.DockerDefaults{
		"Host": conf.RancherInternal("tool-deps", "redis"),
		"Port": 6379,
	}
}

func (r Redis) MarshalDefaults(v interface{}) {
	if rd, ok := v.(*Redis); ok {
		if rd.Name == "" {
			rd.Name = os.Getenv("PROJECT_NAME")
		}
		if rd.Protocol == "" {
			rd.Protocol = "tcp"
		}
		if rd.Port == 0 {
			rd.Port = 6379
		}
		if rd.Password == "" {
			rd.Password = "redis"
		}
		if rd.ConnectTimeout == 0 {
			rd.ConnectTimeout = 10 * time.Second
		}
		if rd.ReadTimeout == 0 {
			rd.ReadTimeout = 10 * time.Second
		}
		if rd.WriteTimeout == 0 {
			rd.WriteTimeout = 10 * time.Second
		}
		if rd.IdleTimeout == 0 {
			rd.IdleTimeout = 240 * time.Second
		}
		if rd.MaxActive == 0 {
			rd.MaxActive = 5
		}
		if rd.MaxIdle == 0 {
			rd.MaxIdle = 3
		}
		if !rd.Wait {
			rd.Wait = true
		}
		if rd.DB == "" {
			rd.DB = "10"
		}
	}
}

func (r *Redis) Prefix(key string) string {
	return fmt.Sprintf("%s::%s::%s", env.GetRuntimeEnv(), r.Name, key)
}

func (r *Redis) TopicFor(topic string, key string) string {
	return r.Prefix(fmt.Sprintf("%s::%s", topic, key))
}

func (r *Redis) Init() {
	if r.cache == nil {
		r.cache = NewRedisCache()
		err := r.cache.Start(r)
		if err != nil {
			panic(err)
		}
	}
}

func (r *Redis) GetCache() *RedisCache {
	return r.cache
}
