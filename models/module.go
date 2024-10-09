package models

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aungmyozaw92/go-graphql/config"
	"github.com/aungmyozaw92/go-graphql/utils"
)

type Module struct {
	ID         int       `gorm:"primary_key" json:"id"`
	Name       string    `gorm:"index;size:100;not null" json:"name" binding:"required"`
	Actions    string    `gorm:"not null" json:"action" binding:"required"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type NewModule struct {
	Name       string `json:"name" binding:"required"`
	Actions    string `json:"action" binding:"required"`
}


//  get ids of roles related to this module / have access
func (module *Module) getRelatedRoleIds(ctx context.Context) ([]int, error) {
	// cache???
	var roleIds []int
	db := config.GetDB()

	err := db.WithContext(ctx).Model(&RoleModule{}).Select("role_id").
		Where("module_id = ?", module.ID).Scan(&roleIds).Error
	if err != nil {
		return nil, err
	}
	return roleIds, nil
}

func (input *NewModule) validate(ctx context.Context, id int) error {
	// name
	if err := utils.ValidateUnique[Module](ctx, "name", input.Name, id); err != nil {
		return err
	}
	return nil
}

func CreateModule(ctx context.Context, input *NewModule) (*Module, error) {

	// ONLY ADMIN can access
	db := config.GetDB()

	// validate module name
	if err := input.validate(ctx, 0); err != nil {
		return nil, err
	}

	module := Module{
		Name:       input.Name,
		Actions:    input.Actions,
	}

	// create module
	tx := db.Begin()
	err := tx.WithContext(ctx).Create(&module).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// remove Cache for Module in Redis 
	if err := utils.RemoveRedisList[Module](); err != nil {
		return nil, err
	}

	return &module, tx.Commit().Error
}

func UpdateModule(ctx context.Context, id int, input *NewModule) (*Module, error) {

	// only admin can access
	db := config.GetDB()
	// check exists
	var count int64
	if err := db.WithContext(ctx).Model(&Module{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return nil, err
	}
	if count <= 0 {
		typeName := utils.GetTypeName[Module]()
		return nil, fmt.Errorf("%s record not found", typeName)
	}

	if err := input.validate(ctx, id); err != nil {
		return nil, err
	}

	module := Module{
		ID:      id,
		Name:    input.Name,
		Actions: input.Actions,
	}

	// update the module
	tx := db.Begin()
	err := tx.WithContext(ctx).Model(&module).Updates(map[string]interface{}{
		"Name":    input.Name,
		"Actions": input.Actions,
	}).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// remove Cache for Module in Redis 
	if err := RemoveRedisBoth(module); err != nil {
		return nil, err
	}

	return &module, tx.Commit().Error
}

func DeleteModule(ctx context.Context, id int) (*Module, error) {

	// only admin can access
	db := config.GetDB()
	var result Module

	err := db.WithContext(ctx).First(&result, id).Error
	if err != nil {
		typeName := utils.GetTypeName[Module]()
		return nil, fmt.Errorf("%s record not found", typeName)
	}

	// delete module
	tx := db.Begin()
	err = tx.WithContext(ctx).Delete(&result).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Delete role module
	err = tx.WithContext(ctx).Where("module_id = ?", id).Delete(&RoleModule{}).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// remove Cache for Module in Redis 
	if err := utils.RemoveRedisList[Module](); err != nil {
		return nil, err
	}

	return &result, tx.Commit().Error
}

func GetModule(ctx context.Context, id int) (*Module, error) {

	// only admin can access
	db := config.GetDB()
	var result Module

	exists, err := utils.GetRedis[Module](id)
	if err != nil {
		return nil, err
	}
	if exists != nil{
		// Module found in cache, return it
		return &result, nil
	}

	if err = db.WithContext(ctx).First(&result, id).Error; err != nil {
		return nil, errors.New("module not found")
	}

	// Cache the module in Redis for future requests
	if err := utils.StoreRedis[Module](result, id); err != nil {
		return nil, err
	}

	return &result, nil
}

func GetModules(ctx context.Context, name *string) ([]*Module, error) {
	db := config.GetDB()
	var results []*Module

	results, err := utils.GetRedisList[Module]()
	if err != nil {
		return nil, err
	}
	if results != nil{
		// Module found in cache, return it
		return results, nil
	}

	if err := db.WithContext(ctx).Find(&results).Error; err != nil {
		return results, errors.New("no modules")
	}

	if err = utils.StoreRedisList[Module](results); err != nil {
		return nil, err
	}

	return results, nil
}
