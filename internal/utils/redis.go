package utils

import (
	"fmt"
	"strconv"

	"github.com/gomodule/redigo/redis"
	_ "github.com/joho/godotenv/autoload"
)

var IataListName = "iata"

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

func ZRANGEBYSCOREWITHSCORES(key string, rangeMin string, rangeMax string) ([]string, error) {

	conn := Pool.Get()
	defer conn.Close()

	data, err := redis.Strings(conn.Do("ZRANGEBYSCORE", key, rangeMin, rangeMax, "WITHSCORES"))
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

/*
To create a Redis lists or to add an nonexistant element to it.
Redis lists are linked lists of strings, sorted by insertion order.
-ex: redis.SetAdd("iatas", []string{"NYC", "NTE"})
*/
func SetAdd(key string, value []byte) error {
	conn := Pool.Get()
	defer conn.Close()

	_, err := redis.Strings(conn.Do("SADD", key, value))
	if err != nil {
		return fmt.Errorf("error setting key %s with value %v", key, value)
	}
	return err
}

/*
-return: if member is a member of the set stored at key.
*/
func SisMember(key string, value byte) (int, error) {
	conn := Pool.Get()
	defer conn.Close()

	data, err := redis.Strings(conn.Do("SISMEMBER", key, value))
	count, _ := strconv.ParseInt(data[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("error evaluating key %s includes %v", key, value)
	}
	return int(count), nil
}

/*
-return: all the members of the set value stored at key.
*/
func SMembers(key string) ([]string, error) {
	conn := Pool.Get()
	defer conn.Close()

	data, err := redis.Strings(conn.Do("SMEMBERS", key))
	if err != nil {
		return nil, fmt.Errorf("error listing elements from key %s", key)
	}
	return data, nil
}

/*
To create a Redis Set.
They are unordered collections of strings.
-ex: redis.SSet("iata:NYC", map[string]int{"temperature": 1, "atmospheric_pressure": 0, "wind_speed": 0})
*/
func HMSet(key string, value []byte) error {
	conn := Pool.Get()
	defer conn.Close()

	_, err := conn.Do("HMSET", key, value)
	if err != nil {
		return fmt.Errorf("error setting key %s with value %s", key, value)
	}
	return err

}

/*
To update the value of a field in a redis hash retrievied by its key.
-ex: redis.HSet("user:1000" "password" 12345)
*/
func HSet(key string, id string, value []byte) error {
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
