package utils

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
	_ "github.com/joho/godotenv/autoload"
)

func Ping() error {

	conn := Pool.Get()
	defer conn.Close()

	_, err := redis.String(conn.Do("PING"))
	if err != nil {
		return fmt.Errorf("cannot 'PING' db: %v", err)
	}
	return nil
}

func Get(key string) ([]byte, error) {

	conn := Pool.Get()
	defer conn.Close()

	var data []byte
	data, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return data, fmt.Errorf("error getting key %s: %v", key, err)
	}
	return data, err
}

func HGet(key string, field string) (interface{}, error) {

	conn := Pool.Get()
	defer conn.Close()

	var data interface{}
	data, err := conn.Do("HGET", key, field)
	fmt.Printf("key: %s\nfield: %s\ndata: %s", key, field, data)
	if err != nil {
		return data, fmt.Errorf("error getting key %s: %v", key, err)
	}
	return data, err
}

func HGetAll(key string) ([]string, error) {

	conn := Pool.Get()
	defer conn.Close()

	data, err := redis.Strings(conn.Do("HGETALL", key))
	fmt.Printf("key: %s\ndata: %s\n", key, data)
	if err != nil {
		return data, fmt.Errorf("error getting key %s: %v", key, err)
	}
	return data, err
}

func ZRANGEBYSCORE(key string, rangeMin string, rangeMax string) ([]string, error) {

	conn := Pool.Get()
	defer conn.Close()

	data, err := redis.Strings(conn.Do("ZRANGEBYSCORE", key, rangeMin, rangeMax))
	fmt.Printf("key: %s\nrangeMin: %s\nrangeMax: %s\n", key, rangeMin, rangeMax)
	if err != nil {
		return nil, fmt.Errorf("error getting key %s: %v", key, err)
	}
	return data, err
}

func Set(key string, value []byte) error {

	conn := Pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, value)
	if err != nil {
		v := string(value)
		if len(v) > 15 {
			v = v[0:12] + "..."
		}
		return fmt.Errorf("error setting key %s to %s: %v", key, v, err)
	}
	return err
}

func ZSet(key string, field string, value string) error {
	conn := Pool.Get()
	defer conn.Close()

	_, err := conn.Do("ZADD", key, field, value)
	if err != nil {
		return fmt.Errorf("error setting key %s with value %s", key, value)
	}
	return err
}

func HSET(key string, id string, value []byte) error {
	conn := Pool.Get()
	defer conn.Close()

	_, err := conn.Do("HSET", key, id, value)
	if err != nil {
		return fmt.Errorf("error setting key %s and id %s with value %s", key, id, value)
	}
	return err
}

func Exists(key string) (bool, error) {

	conn := Pool.Get()
	defer conn.Close()

	ok, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return ok, fmt.Errorf("error checking if key %s exists: %v", key, err)
	}
	return ok, err
}

func Delete(key string) error {

	conn := Pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	return err
}

func GetKeys(pattern string) ([]string, error) {

	conn := Pool.Get()
	defer conn.Close()

	iter := 0
	keys := []string{}
	for {
		arr, err := redis.Values(conn.Do("SCAN", iter, "MATCH", pattern))
		if err != nil {
			return keys, fmt.Errorf("error retrieving '%s' keys", pattern)
		}

		iter, _ = redis.Int(arr[0], nil)
		k, _ := redis.Strings(arr[1], nil)
		keys = append(keys, k...)

		if iter == 0 {
			break
		}
	}

	return keys, nil
}

func Incr(counterKey string) (int, error) {

	conn := Pool.Get()
	defer conn.Close()

	return redis.Int(conn.Do("INCR", counterKey))
}
