package config

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/bsm/redislock"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var (
	rdb    *redis.Client
	locker *redislock.Client
)
var ctx = context.Background()

func GetRedisDB() *redis.Client {
	return rdb
}

func GetRedisLock() *redislock.Client {
	return locker
}

func GetRedisContext() context.Context {
	return ctx
}

func GetRedisObject(key string, dest interface{}) (bool, error) {
	// fmt.Printf("	(Redis) Getting object of `%s`\n", key)
	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return false, nil
		}
		return false, err
	}
	err = json.Unmarshal([]byte(val), &dest)
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetRedisValue(key string) (string, bool, error) {
	// fmt.Printf("	(Redis) Getting value of `%s`\n", key)
	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return "", false, nil
		}
		return "", false, err
	}
	return val, true, nil
}

func SetRedisObject(key string, obj interface{}, exp time.Duration) error {
	// fmt.Printf("	(Redis) Setting object `%s`:%+v\n", key, obj)
	objInByte, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	if err = rdb.Set(ctx, key, objInByte, exp).Err(); err != nil {
		return err
	}
	return nil
}

// store key in a set for faster adding & retrieving
func AddRedisSet(setKey string, member string) error {
	if err := rdb.SAdd(ctx, setKey, member).Err(); err != nil {
		return err
	}
	return nil
}

func GetRedisSetMembers(setKey string) ([]string, error) {
	return rdb.SMembers(ctx, setKey).Result()
}

func RemoveRedisSetMember(setKey string, member string) error {
	return rdb.SRem(ctx, setKey, member).Err()
}

func SetRedisValue(key string, value string, exp time.Duration) error {
	// fmt.Printf("	(Redis) Setting value `%s`:%s\n", key, value)
	// rdb.Set
	return rdb.Set(ctx, key, value, exp).Err()
}

func RemoveRedisKey(keys ...string) error {
	// fmt.Printf("	(Redis) Removing `%v`\n", keys)
	_, err := rdb.Del(ctx, keys...).Result()
	return err
}

func ClearRedis(ctx context.Context) error {
	cmd := rdb.FlushAll(ctx)
	return cmd.Err()
}

// add one and returns it, while storing the updated value
func GetRedisCounter(ctx context.Context, key string) (int64, error) {
	// exists, err := GetRedisObject(key, &value)
	// if err != nil {
	// 	return 0, err
	// }
	// if !exists {
	// 	// TODO: get from database
	// 	return 0, nil
	// }
	// result, err := rdb.Incr(ctx, key).Result()
	// if err != nil {
	// 	return 0, err
	// }

	return rdb.Incr(ctx, key).Result()
}

func init() {
	// Load env from .env
	godotenv.Load()
	connectRedis()
	locker = redislock.New(rdb)
}

func connectRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: "",
		DB:       1, // use default DB
		PoolSize: 100,
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		panic("Fail to Connect Redis")
	}
}
