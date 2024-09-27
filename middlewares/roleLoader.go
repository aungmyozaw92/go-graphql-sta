package middlewares

import (
	"context"

	"github.com/aungmyozaw92/go-graphql/models"
	"github.com/graph-gophers/dataloader/v7"
	"gorm.io/gorm"
)

// roleReader reads Roles from a database
type roleReader struct {
	db *gorm.DB
}

// getRoles implements a batch function that can retrieve many roles by ID,
// for use in a dataloader

func (u *roleReader) getRoles(ctx context.Context, ids []int) []*dataloader.Result[*models.Role] {
	var results []*models.Role

	err := u.db.WithContext(ctx).Where("id IN ?", ids).Find(&results).Error
	if err != nil {
		// Instead of returning []error, create a single error for the dataloader.Result
		return handleError[*models.Role](len(ids), err)
	}

	loaderResults := make([]*dataloader.Result[*models.Role], 0, len(ids))
	for _, id := range ids {
		if id == 0 {
			loaderResults = append(loaderResults, &dataloader.Result[*models.Role]{Data: &models.Role{}})
		} else {
			for _, result := range results {
				if result.ID == id {
					loaderResults = append(loaderResults, &dataloader.Result[*models.Role]{Data: result})
					break
				}
			}
		}
	}
	return loaderResults
}

// GetRole returns single Role by id efficiently

func GetRole(ctx context.Context, id int) (*models.Role, error) {
	loaders := For(ctx)
	return loaders.RoleLoader.Load(ctx, id)()
}

// GetRoles returns many Roles by ids efficiently
func GetRoles(ctx context.Context, ids []int) ([]*models.Role, []error) {
	loaders := For(ctx)
	return loaders.RoleLoader.LoadMany(ctx, ids)()
}
