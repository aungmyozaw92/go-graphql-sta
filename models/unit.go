package models

import (
	"context"
	"errors"
	"time"

	"github.com/aungmyozaw92/go-graphql/config"
	"github.com/aungmyozaw92/go-graphql/utils"
)

type Unit struct {
	ID           int       `gorm:"primary_key" json:"id"`
	Name         string    `gorm:"size:20;not null" json:"name" binding:"required"`
	Abbreviation string    `gorm:"size:10;not null" json:"abbreviation" binding:"required"`
	Precision    Precision `gorm:"type:enum('0','1','2','3','4');default:'0';size:1;not null" json:"precision" binding:"required"`
	IsActive     *bool     `gorm:"not null;default:true" json:"is_active"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type NewUnit struct {
	Name         string    `json:"name" binding:"required"`
	Abbreviation string    `json:"abbreviation" binding:"required"`
	Precision    Precision `json:"precision" binding:"required"`
}

type UnitsEdge Edge[Unit]
type UnitsConnection struct {
	PageInfo *PageInfo      `json:"pageInfo"`
	Edges    []*UnitsEdge 	`json:"edges"`
}

// node
// returns decoded curosr string
func (pu Unit) GetCursor() string {
	return pu.CreatedAt.String()
}

func (input *NewUnit) validate(ctx context.Context, id int) error {
	if err := utils.ValidateUnique[Unit](ctx, "name", input.Name, id); err != nil {
		return err
	}
	if err := utils.ValidateUnique[Unit](ctx, "abbreviation", input.Abbreviation, id); err != nil {
		return err
	}

	return nil
}

func CreateUnit(ctx context.Context, input *NewUnit) (*Unit, error) {

	if err := input.validate(ctx, 0); err != nil {
		return nil, err
	}

	unit := Unit{
		Name:         input.Name,
		Abbreviation: input.Abbreviation,
		Precision:    input.Precision,
	}

	db := config.GetDB()

	err := db.WithContext(ctx).Create(&unit).Error

	if err != nil {
		return nil, err
	}

	// remove Cache for Unit in Redis 
	if err := utils.RemoveRedisList[Unit](); err != nil {
		return nil, err
	}

	return &unit, nil
}

func UpdateUnit(ctx context.Context, id int, input *NewUnit) (*Unit, error) {

	if err := input.validate(ctx, id); err != nil {
		return nil, err
	}

	var unit Unit

	db := config.GetDB()
	err := db.WithContext(ctx).First(&unit, id).Error

	if err != nil {
		return nil, err
	}
	
	err = db.WithContext(ctx).Model(&unit).Updates(map[string]interface{}{
		"Name":         input.Name,
		"Abbreviation": input.Abbreviation,
		"Precision":    input.Precision,
	}).Error
	if err != nil {
		return nil, err
	}


	// remove Cache for Module in Redis 
	if err := RemoveRedisBoth(unit); err != nil {
		return nil, err
	}

	return &unit, nil
}

func DeleteUnit(ctx context.Context, id int) (*Unit, error) {

	
	var unit Unit

	db := config.GetDB()
	err := db.WithContext(ctx).First(&unit, id).Error
	if err != nil {
		return nil, err
	}

	// don't delete if Unit is used by  or  variant
	count, err := utils.ResourceCountWhere[Product](ctx, "unit_id = ?", id)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("used by product")
	}
	// count, err = utils.ResourceCountWhere[ProductVariant](ctx, "unit_id = ?", id)
	// if err != nil {
	// 	return nil, err
	// }
	// if count > 0 {
	// 	return nil, errors.New("used by product variant")
	// }

	// db action
	err = db.WithContext(ctx).Delete(&unit).Error
	if err != nil {
		return nil, err
	}

	// remove Cache for Unit in Redis 
	if err := utils.RemoveRedisList[Unit](); err != nil {
		return nil, err
	}

	return &unit, nil
}

func ToggleActiveUnit(ctx context.Context, id int, isActive bool) (*Unit, error) {

	var unit Unit

	db := config.GetDB()
	err := db.WithContext(ctx).First(&unit, id).Error

	if err != nil {
		return nil, err
	}
	
	err = db.WithContext(ctx).Model(&unit).Updates(map[string]interface{}{
		"IsActive":   isActive,
	}).Error
	if err != nil {
		return nil, err
	}

	// remove Cache for Unit in Redis 
	if err := RemoveRedisBoth(unit); err != nil {
		return nil, err
	}

	return &unit, nil
}


func GetUnit(ctx context.Context, id int) (*Unit, error) {

	return GetResource[Unit](ctx, id)
}

func GetUnits(ctx context.Context, name *string) ([]*Unit, error) {

	results, err := GetResources[Unit](ctx, "name")

	if err != nil {
		return nil, err
	}

	return results, nil
}


func PaginateUnit(ctx context.Context, limit *int, after *string, name *string) (*UnitsConnection, error) {

	db := config.GetDB()
	dbCtx := db.WithContext(ctx)

	if name != nil && *name != "" {
		dbCtx.Where("name LIKE ?", "%"+*name+"%")
	}

	// edges, pageInfo, err := FetchPagePureCursor[Unit](dbCtx, *limit, after, "name", ">")
	edges, pageInfo, err := FetchPageCompositeCursor[Unit](dbCtx, *limit, after, "name", ">")
	if err != nil {
		return nil, err
	}

	var UnitsConnection UnitsConnection
	UnitsConnection.PageInfo = pageInfo

	for _, edge := range edges {
		UnitEdge := UnitsEdge(edge)
		UnitsConnection.Edges = append(UnitsConnection.Edges, &UnitEdge)
	}

	return &UnitsConnection, err
}