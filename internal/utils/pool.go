package utils

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gomodule/redigo/redis"
)

/*
A connection pool is a connection server!
When such a server is started, it opens a number of connections to the database server.
The clients address the pool, which assigns them a connection.
*/
var (
	Pool *redis.Pool
)

/*
Initialized redis pool is stored in Pool variable.
Pool connects by default localhost:6379.
This can be changed by passing REDIS_HOST environment variable.
*/
func init() {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = ":6379"
	}
	Pool = newPool(redisHost)
	cleanupHook()
}

func newPool(server string) *redis.Pool {

	return &redis.Pool{

		//This is the max number of connections that can be idle in the pool without being immediately evicted (closed).
		MaxIdle: 3,
		//Idle connections are closed after they are idle for the IdleTimeout duration.
		IdleTimeout: 240 * time.Second,
		//Dial connects to the Redis server at the given network and address using the specified options.
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		//Controls whether or not the connection is tested before it is returned from the pool.
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

/*
Shut down the pool when it receives the following notifications
+ a SIGTERM
+ a SIGKILL
*/
func cleanupHook() {
	//create a channel to receive these notifications
	c := make(chan os.Signal, 1)
	//registers the given channel to receive notifications of the specified signals.
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	//wait until it gets the expected signal
	go func() {
		<-c
		Pool.Close()
		os.Exit(0)
	}()
}
