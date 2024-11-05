package middlewares

import (
	"context"

	"github.com/aungmyozaw92/go-graphql/models"
	"github.com/graph-gophers/dataloader/v7"
	"gorm.io/gorm"
)

// unitReader reads Unit from a database
type unitReader struct {
	db *gorm.DB
}

// getUnits implements a batch function that can retrieve many Units by ID,
// for use in a dataloader

func (u *unitReader) getUnits(ctx context.Context, ids []int) []*dataloader.Result[*models.Unit] {
	var results []*models.Unit

	err := u.db.WithContext(ctx).Where("id IN ?", ids).Find(&results).Error
	if err != nil {
		// Instead of returning []error, create a single error for the dataloader.Result
		return handleError[*models.Unit](len(ids), err)
	}

	loaderResults := make([]*dataloader.Result[*models.Unit], 0, len(ids))
	for _, id := range ids {
		if id == 0 {
			loaderResults = append(loaderResults, &dataloader.Result[*models.Unit]{Data: &models.Unit{}})
		} else {
			for _, result := range results {
				if result.ID == id {
					loaderResults = append(loaderResults, &dataloader.Result[*models.Unit]{Data: result})
					break
				}
			}
		}
	}
	return loaderResults
}

// GetUnit returns single Unit by id efficiently

func GetUnit(ctx context.Context, id int) (*models.Unit, error) {
	loaders := For(ctx)
	return loaders.UnitLoader.Load(ctx, id)()
}

// GetUnits returns many Units by ids efficiently
func GetUnits(ctx context.Context, ids []int) ([]*models.Unit, []error) {
	loaders := For(ctx)
	return loaders.UnitLoader.LoadMany(ctx, ids)()
}
