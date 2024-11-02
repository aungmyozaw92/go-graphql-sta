package models

import (
	"context"
	"errors"
	"time"

	"github.com/aungmyozaw92/go-graphql/config"
	"github.com/aungmyozaw92/go-graphql/utils"
	"gorm.io/gorm"
)

type Category struct {
	ID               int       `gorm:"primary_key" json:"id"`
	Name             string    `gorm:"index;size:100;not null" json:"name" binding:"required"`
	ParentCategoryId int       `gorm:"index;not null" json:"parentCategoryId"`
	IsActive         *bool     `gorm:"not null;default:true" json:"is_active"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type NewCategory struct {
	Name             string `json:"name" binding:"required"`
	ParentCategoryId int    `json:"parentCategoryId" binding:"required"`
}

// implements methods for pagination
type CategoriesEdge Edge[Category]
type CategoriesConnection struct {
	PageInfo *PageInfo                `json:"pageInfo"`
	Edges    []*CategoriesEdge 		  `json:"edges"`
}

// node
// returns decoded curosr string
func (pc Category) GetCursor() string {
	return pc.CreatedAt.String()
}

// validate input for both create & update. (id = 0 for create)

func (input *NewCategory) validate(ctx context.Context, id int) error {
	if id > 0 {
		if id == input.ParentCategoryId {
			return errors.New("self-parent not allowed")
		}
	}
	// name
	if err := utils.ValidateUnique[Category](ctx, "name", input.Name, id); err != nil {
		return err
	}
	// parent category
	if input.ParentCategoryId > 0 {
		if err := utils.ValidateResourceId[Category](ctx, input.ParentCategoryId); err != nil {
			return errors.New("parent not found")
		}
	}
	return nil
}

func CreateCategory(ctx context.Context, input *NewCategory) (*Category, error) {

	if err := input.validate(ctx, 0); err != nil {
		return nil, err
	}

	category := Category{
		Name:             input.Name,
		ParentCategoryId: input.ParentCategoryId,
		IsActive:         utils.NewTrue(),
	}

	db := config.GetDB()
	err := db.WithContext(ctx).Create(&category).Error
	if err != nil {
		return nil, err
	}

	// remove Cache for Category in Redis 
	if err := utils.RemoveRedisList[Category](); err != nil {
		return nil, err
	}

	return &category, nil
}


func UpdateCategory(ctx context.Context, id int, input *NewCategory) (*Category, error) {

	var category Category
	if err := input.validate(ctx, id); err != nil {
		return nil, err
	}

	db := config.GetDB()

	err := db.WithContext(ctx).First(&category, id).Error
	if err != nil {
		return nil, err
	}

	err = db.WithContext(ctx).Model(&category).Updates(map[string]interface{}{
		"Name":             input.Name,
		"ParentCategoryId": input.ParentCategoryId,
	}).Error

	if err != nil {
		return nil, err
	}

	// remove Cache for Module in Redis 
	if err := RemoveRedisBoth(category); err != nil {
		return nil, err
	}

	return &category, nil
}

func DeleteCategory(ctx context.Context, id int) (*Category, error) {

	var category Category

	db := config.GetDB()

	err := db.WithContext(ctx).First(&category, id).Error
	if err != nil {
		return nil, err
	}

	// don't delete if Category has childern
	count, err := utils.ResourceCountWhere[Category](ctx, "parent_category_id = ?", id)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("category has children")
	}

	err = db.WithContext(ctx).Delete(&category).Error
	if err != nil {
		return nil, err
	}

	// remove Cache for Module in Redis 
	if err := RemoveRedisBoth(category); err != nil {
		return nil, err
	}

	return &category, nil
}

func ToggleActiveCategory(ctx context.Context, id int, isActive bool) (*Category, error) {

	var category Category
	db := config.GetDB()
	if err := db.WithContext(ctx).Where("id = ?", id).Find(&category).Error; err != nil {
		return nil, errors.New("category not found")
	}

	tx := db.Begin()
	if err := tx.Model(&category).Updates(map[string]interface{}{
		"is_active": isActive,
	}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := toggleChildrenCategories(ctx, tx, id, isActive); err != nil {
		tx.Rollback()
		return &category, err
	}

	// remove Cache for Module in Redis 
	if err := RemoveRedisBoth(category); err != nil {
		return nil, err
	}

	return &category, tx.Commit().Error
}


// toggle children of the parent recursively, parent is assumed to have toggled
func toggleChildrenCategories(ctx context.Context, tx *gorm.DB, parentId int, isActive bool) error {
	// get children ids
	// toggle them
	// toggle children of each child
	// break when a parent has no children

	var childrenIds []int
	if err := tx.WithContext(ctx).
		Model(&Category{}).
		Where("parent_category_id = ?", parentId).
		Select("id").
		Scan(&childrenIds).Error; err != nil {
		return err
	}

	// base case
	// break when parent has no children
	if len(childrenIds) == 0 {
		return nil
	}

	if err := tx.WithContext(ctx).Model(&Category{}).
		Where("id IN ?", childrenIds).Updates(map[string]interface{}{
		"is_active": isActive,
	}).Error; err != nil {
		return err
	}

	for _, childId := range childrenIds {
		// each child becomes a parent
		if err := toggleChildrenCategories(ctx, tx, childId, isActive); err != nil {
			return err
		}
	}
	return nil
}


func GetCategory(ctx context.Context, id int) (*Category, error) {

	return GetResource[Category](ctx, id)
}

func GetCategories(ctx context.Context, name *string) ([]*Category, error) {

	results, err := GetResources[Category](ctx, "name")

	if err != nil {
		return nil, err
	}

	return results, nil
}

func PaginateCategory(ctx context.Context, limit *int, after *string, name *string, parentCategoryId *int) (*CategoriesConnection, error) {

	db := config.GetDB()
	dbCtx := db.WithContext(ctx)

	
	if name != nil && *name != "" {
		dbCtx.Where("name LIKE ?", "%"+*name+"%")
	}
	if parentCategoryId != nil && *parentCategoryId > 0 {
		dbCtx.Where("parent_category_id = ?", *parentCategoryId)
	}

	edges, pageInfo, err := FetchPageCompositeCursor[Category](dbCtx, *limit, after, "created_at", "<")

	if err != nil {
		return nil, err
	}
	var categoriesConnection CategoriesConnection
	categoriesConnection.PageInfo = pageInfo
	for _, edge := range edges {
		categoryEdge := CategoriesEdge(edge)
		categoriesConnection.Edges = append(categoriesConnection.Edges, &categoryEdge)
	}
	return &categoriesConnection, nil
}