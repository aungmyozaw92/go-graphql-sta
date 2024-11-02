package models

import "github.com/aungmyozaw92/go-graphql/utils"

type RedisCleaner interface {
	RemoveInstanceRedis() error // remove one
	RemoveAllRedis() error      // remove list & map if exists
}

// remove both item & list + map
func RemoveRedisBoth[T RedisCleaner](obj T) error {
	if err := obj.RemoveInstanceRedis(); err != nil {
		return err
	}
	if err := obj.RemoveAllRedis(); err != nil {
		return err
	}
	return nil
}

func (obj Module) RemoveInstanceRedis() error {
	if err := utils.RemoveRedisItem[Module](obj.ID); err != nil {
		return err
	}
	return nil
}

func (obj Module) RemoveAllRedis() error {
	if err := utils.RemoveRedisList[Module](); err != nil {
		return err
	}
	return nil
}

func (obj Unit) RemoveInstanceRedis() error {
	if err := utils.RemoveRedisItem[Unit](obj.ID); err != nil {
		return err
	}
	return nil
}

func (obj Unit) RemoveAllRedis() error {
	if err := utils.RemoveRedisList[Unit](); err != nil {
		return err
	}
	return nil
}

func (obj Category) RemoveInstanceRedis() error {
	if err := utils.RemoveRedisItem[Category](obj.ID); err != nil {
		return err
	}
	return nil
}

func (obj Category) RemoveAllRedis() error {
	if err := utils.RemoveRedisList[Category](); err != nil {
		return err
	}
	return nil
}