package utils

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/aungmyozaw92/go-graphql/config"
)

func GetCacheLifespan() time.Duration {
	lifespan, err := strconv.Atoi(os.Getenv("CACHE_LIFESPAN"))
	if err != nil {
		lifespan = 1
	}
	return time.Duration(lifespan) * time.Hour
}

/* generic functions */

func GetTypeName[T any]() string {
	var v T
	typeOfT := reflect.TypeOf(v)
	return typeOfT.Name()
}

// get type name of struct
func GetType(i interface{}) string {
	return reflect.TypeOf(i).Name()
}


// store instance, obj should be a pointer
func StoreRedis[T any](obj any, id int) error {
	typeName := GetTypeName[T]()
	key := typeName + ":" + fmt.Sprint(id)

	return config.SetRedisObject(key, &obj, GetCacheLifespan())
}

// store object
func StoreRedisList[T any](obj any) error {
	var key string
	key = GetTypeName[T]() + "List"
	return config.SetRedisObject(key, &obj,  GetCacheLifespan())
}

// get from redis
// returns nil if does not exist
func GetRedis[T any](id int) (*T, error) {
	var result *T
	key := GetTypeName[T]() + ":" + fmt.Sprint(id)
	exists, err := config.GetRedisObject(key, &result)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, nil
	}
	return result, nil
}

// retrieve a list.
// businessId can be empty
func GetRedisList[T any]() ([]*T, error) {
	var key string
		key = GetTypeName[T]() + "List"
	

	var result []*T
	exists, err := config.GetRedisObject(key, &result)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, nil
	}
	return result, nil
}

// remove an instance, Type:$id
func RemoveRedisItem[T any](id int) error {
	key := GetTypeName[T]() + ":" + fmt.Sprint(id)
	return config.RemoveRedisKey(key)
}

// clear list, TypeList
func RemoveRedisList[T any]() error {
	var key string = GetTypeName[T]() + "List"
	return config.RemoveRedisKey(key)
}

