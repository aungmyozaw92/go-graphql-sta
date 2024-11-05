package middlewares

import (
	"context"

	"github.com/aungmyozaw92/go-graphql/models"
	"github.com/graph-gophers/dataloader/v7"
	"gorm.io/gorm"
)

// categoryReader reads Categories from a database
type categoryReader struct {
	db *gorm.DB
}

// getCategories implements a batch function that can retrieve many Categories by ID,
// for use in a dataloader

func (u *categoryReader) getCategories(ctx context.Context, ids []int) []*dataloader.Result[*models.Category] {
	var results []*models.Category

	err := u.db.WithContext(ctx).Where("id IN ?", ids).Find(&results).Error
	if err != nil {
		// Instead of returning []error, create a single error for the dataloader.Result
		return handleError[*models.Category](len(ids), err)
	}

	loaderResults := make([]*dataloader.Result[*models.Category], 0, len(ids))
	for _, id := range ids {
		if id == 0 {
			loaderResults = append(loaderResults, &dataloader.Result[*models.Category]{Data: &models.Category{}})
		} else {
			for _, result := range results {
				if result.ID == id {
					loaderResults = append(loaderResults, &dataloader.Result[*models.Category]{Data: result})
					break
				}
			}
		}
	}
	return loaderResults
}

// GetCategory returns single Category by id efficiently

func GetCategory(ctx context.Context, id int) (*models.Category, error) {
	loaders := For(ctx)
	return loaders.CategoryLoader.Load(ctx, id)()
}

// GetCategories returns many Categories by ids efficiently
func GetCategories(ctx context.Context, ids []int) ([]*models.Category, []error) {
	loaders := For(ctx)
	return loaders.CategoryLoader.LoadMany(ctx, ids)()
}
