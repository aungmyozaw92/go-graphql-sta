package models

import (
	"context"
	"errors"
	"time"

	"github.com/aungmyozaw92/go-graphql/config"
	"github.com/aungmyozaw92/go-graphql/utils"
	"github.com/shopspring/decimal"
)


type Product struct {
	ID                  int               `gorm:"primary_key" json:"id"`
	BusinessId          string            `gorm:"index;not null" json:"business_id" binding:"required"`
	Name                string            `gorm:"size:100;not null" json:"name" binding:"required"`
	Description         string            `gorm:"type:text" json:"description"`
	CategoryId          int               `gorm:"index;not null;default:0" json:"category_id"`
	Images              []*Image          `gorm:"polymorphic:Reference" json:"-"`
	UnitId              int               `json:"product_unit_id"`
	SupplierId          int               `json:"supplier_id"`
	Sku                 string            `gorm:"size:100;not null" json:"sku" binding:"required"`
	Barcode             string            `gorm:"index;size:100;not null" json:"barcode" binding:"required"`
	SalesPrice          decimal.Decimal   `gorm:"type:decimal(20,4);default:0" json:"sales_price"`
	PurchasePrice       decimal.Decimal   `gorm:"type:decimal(20,4);default:0" json:"purchase_price"`
	IsActive            *bool             `gorm:"not null;default:true" json:"is_active"`
	IsBatchTracking     *bool             `gorm:"not null;default:false" json:"is_batch_traking"`
	CreatedAt 			time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt 			time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type NewProduct struct {
	Name                string                   `json:"name" binding:"required"`
	Description         string                   `json:"description"`
	CategoryId          int                      `json:"category_id"`
	Images              []*NewImage              `json:"image_urls"`
	UnitId              int                      `json:"product_unit_id"`
	SupplierId          int                      `json:"supplier_id"`
	Sku                 string                   `json:"sku" binding:"required"`
	Barcode             string                   `json:"barcode" binding:"required"`
	SalesPrice          decimal.Decimal          `json:"sales_price"`
	PurchasePrice       decimal.Decimal          `json:"purchase_price"`
	IsBatchTracking     *bool                    `json:"is_batch_traking"`
}

type ProductsEdge Edge[Product]

type ProductsConnection struct {
	PageInfo *PageInfo
	Edges []*ProductsEdge
}

// implements Node
func (p Product) GetCursor() string {
	return p.CreatedAt.String()
}

// validate input for both create & update. (id = 0 for create)

func (input *NewProduct) validate(ctx context.Context, id int) error {
	if err := utils.ValidateUnique[Product](ctx, "name", input.Name, id); err != nil {
		return err
	}
	// exists category
	if input.CategoryId > 0 {
		if err := utils.ValidateResourceId[Category](ctx, input.CategoryId); err != nil {
			return errors.New("product category not found")
		}
	}

	// exists unit
	if input.UnitId > 0 {
		if err := utils.ValidateResourceId[Unit](ctx, input.UnitId); err != nil {
			return errors.New("unit not found")
		}
	}

	return nil
}

func CreateProduct(ctx context.Context, input *NewProduct) (*Product, error) {

	// validate product
	if err := input.validate(ctx, 0); err != nil {
		return nil, err
	}

	// construct Images
	images, err := mapNewImages(input.Images, "products", 0)
	if err != nil {
		return nil, err
	}
	// store product
	product := Product{
		Name:                input.Name,
		Description:         input.Description,
		CategoryId:          input.CategoryId,
		UnitId:              input.UnitId,
		SupplierId:          input.SupplierId,
		Sku:                 input.Sku,
		Barcode:             input.Barcode,
		SalesPrice:          input.SalesPrice,
		PurchasePrice:       input.PurchasePrice,
		IsActive:            utils.NewTrue(),
		IsBatchTracking:     input.IsBatchTracking,
		// asssociation
		Images:    images,
	}

	db := config.GetDB()
	tx := db.Begin()

	err = tx.WithContext(ctx).Create(&product).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return &product, nil
}


func DeleteProduct(ctx context.Context, id int) (*Product, error) {
	
	var result Product

	// db action
	db := config.GetDB()
	tx := db.Begin()

	err := tx.WithContext(ctx).Preload("Images").First(&result, id).Error

	if err != nil {
		return nil, err
	}

	for _, img := range result.Images {
		if err := img.Delete(tx, ctx); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// db action
	err = tx.WithContext(ctx).Delete(&result).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &result, nil
}

func GetProduct(ctx context.Context, id int) (*Product, error) {
	return GetResource[Product](ctx, id)
}

func GetProducts(ctx context.Context, name *string) ([]*Product, error) {
	db := config.GetDB()
	var results []*Product


	dbCtx := db.WithContext(ctx)
	if name != nil && len(*name) > 0 {
		dbCtx = dbCtx.Where("name LIKE ?", "%"+*name+"%")
	}

	err := dbCtx.Order("name").
		Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}