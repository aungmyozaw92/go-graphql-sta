package models

import (
	"context"

	"github.com/aungmyozaw92/go-graphql/config"
	"github.com/aungmyozaw92/go-graphql/utils"
)

// (may return RecordNotFound error)
func GetResource[T any](ctx context.Context, id int, associations ...string) (*T, error) {

	// find in redis
	result, err := utils.GetRedis[T](id)
	if err != nil {
		return nil, err
	}
	// if not found in redis
	if result == nil {
		// fetch from db
		// result, err = utils.FetchModel[T](ctx, id, associations...)

		db := config.GetDB()
		dbCtx := db.WithContext(ctx)
		// preloading
		for _, field := range associations {
			dbCtx.Preload(field)
		}
		var result T
		err := dbCtx.First(&result, id).Error
		if err != nil {
			return nil, err
		}

		// store in redis
		if err := utils.StoreRedis[T](result, id); err != nil {
			return nil, err
		}
	} 

	return result, nil
}

// list all resources, redis or db, cache result
func GetResources[Model any](ctx context.Context, orders ...string) ([]*Model, error) {

	// first try redis cache
	results, err := utils.GetRedisList[Model]()
	if err != nil {
		return nil, err
	}
	// if not exists in redis
	if results == nil {
		// fetch from db
		db := config.GetDB()
		var model Model
		dbCtx := db.WithContext(ctx)
		for _, order := range orders {
			dbCtx.Order(order)
		}
		// db query
		if err = dbCtx.Model(&model).Find(&results).Error; err != nil {
			return nil, err
		}

		// caching the result
		if err := utils.StoreRedisList[Model](results); err != nil {
			return nil, err
		}
	}

	return results, nil
}