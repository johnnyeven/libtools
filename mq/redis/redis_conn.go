package redis

import (
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
)

func ConnRedis(pool *redis.Pool) *RedisConn {
	return &RedisConn{
		pool: pool,
	}
}

type RedisConn struct {
	pool     *redis.Pool
	commands []*redisCommand
}

type redisCommand struct {
	cmd  string
	args []interface{}
}

func (c *RedisConn) Do(cmd string, args ...interface{}) (interface{}, error) {
	conn := c.pool.Get()
	defer conn.Close()

	res, err := conn.Do(cmd, args...)
	if err != nil {
		logrus.WithField("error", err.Error()).Println(append([]interface{}{cmd}, args...))
		return nil, err
	}

	return res, err
}

func (c *RedisConn) Send(cmd string, args ...interface{}) {
	c.commands = append(c.commands, &redisCommand{
		cmd:  cmd,
		args: args,
	})
}

func (c *RedisConn) Exec() (interface{}, error) {
	if len(c.commands) == 0 {
		return nil, nil
	}
	conn := c.pool.Get()
	defer conn.Close()

	c.Send("MULTI")

	for _, cmd := range c.commands {
		err := conn.Send(cmd.cmd, cmd.args...)
		if err != nil {
			return nil, err
		}
	}

	res, err := conn.Do("EXEC")
	if err != nil {
		logrus.WithField("error", err.Error()).Print()
		return nil, err
	}

	c.commands = nil

	return res, err
}
