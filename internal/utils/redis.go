package utils

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
	_ "github.com/joho/godotenv/autoload"
)

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