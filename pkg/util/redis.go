package util

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/gomodule/redigo/redis"
)

var pool *redis.Pool
var lock sync.RWMutex

// Redis initialization
func InitRedis() error {
	lock = sync.RWMutex{}
	if err := initPool(); err != nil {
		log.Fatalf("failed initialize db pool. err %v", err)
		return nil
	}
	if os.Getenv("TEST_MODE") == "true" {
		if err := resetDB(); err != nil {
			return err
		}
	}
	return nil
}

func resetDB() error {
	conn := pool.Get()
	defer conn.Close()

	if conn == nil {
		log.Printf("Error: failed connect to redis pool. possibly redis is not setup.")
		return fmt.Errorf("failed to connect redis poll")
	}

	lock.Lock()
	_, err := conn.Do("flushall")
	if err != nil {
		log.Fatalf("failed reset redis db within the environment. err %v", err)
		return err
	}
	lock.Unlock()
	log.Printf("Successfully flush residue data.")
	return nil
}

func initPool() error {
	pool = &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", "localhost:6379")
			if err != nil {
				log.Printf("ERROR: fail init redis: %s", err.Error())
				os.Exit(1)
			}
			return conn, err
		},
	}
	if pool == nil {
		return fmt.Errorf("failed to initialize redis pool connection")
	}
	return nil
}

func Set(key string, value uint64) error {
	if pool == nil {
		log.Printf("redis pool is not initialized yet.")
		return fmt.Errorf("redis pool is not initialized yet")
	}
	// get conn and put back when exit from method
	conn := pool.Get()
	defer conn.Close()

	if conn == nil {
		log.Printf("Error: failed connect to redis pool. possibly redis is not setup.")
		return fmt.Errorf("failed to connect redis poll")
	}

	lock.Lock()
	valueString := strconv.FormatUint(value, 10)
	_, err := conn.Do("SET", key, valueString)
	lock.Unlock()

	if err != nil {
		log.Printf("ERROR: fail set key %s value %v, error %s", key, value, err.Error())
		return err
	}

	log.Printf("[Redis][Set]OK: set key %s value uint64 %v string %s", key, value, valueString)
	return nil
}

func Get(key string) (string, bool) {
	// get conn and put back when exit from method
	if pool == nil {
		log.Printf("redis pool is not initialized yet.")
		return "", false
	}
	conn := pool.Get()
	defer conn.Close()

	if conn == nil {
		log.Printf("failed connect to redis pool. possibly redis is not setup.")
		return "", false
	}

	lock.RLock()
	v, err := redis.String(conn.Do("GET", key))
	lock.RUnlock()
	if err != nil {
		log.Printf("ERROR: fail get key %s, error %s", key, err.Error())
		return "", false
	}

	return v, true
}
