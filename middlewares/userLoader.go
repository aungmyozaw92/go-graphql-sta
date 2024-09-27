package middlewares

import (
	"context"

	"github.com/aungmyozaw92/go-graphql/models"
	"github.com/graph-gophers/dataloader/v7"
	"gorm.io/gorm"
)

// userReader reads Users from a database
type userReader struct {
	db *gorm.DB
}

// getUsers implements a batch function that can retrieve many users by ID,
// for use in a dataloader

func (u *userReader) getUsers(ctx context.Context, ids []int) []*dataloader.Result[*models.User] {
	var results []*models.User

	err := u.db.WithContext(ctx).Where("id IN ?", ids).Find(&results).Error
	if err != nil {
		// Instead of returning []error, create a single error for the dataloader.Result
		return handleError[*models.User](len(ids), err)
	}

	loaderResults := make([]*dataloader.Result[*models.User], 0, len(ids))
	for _, id := range ids {
		if id == 0 {
			loaderResults = append(loaderResults, &dataloader.Result[*models.User]{Data: &models.User{}})
		} else {
			for _, result := range results {
				if result.ID == id {
					loaderResults = append(loaderResults, &dataloader.Result[*models.User]{Data: result})
					break
				}
			}
		}
	}
	return loaderResults
}

// GetUser returns single user by id efficiently

func GetUser(ctx context.Context, id int) (*models.User, error) {
	loaders := For(ctx)
	return loaders.UserLoader.Load(ctx, id)()
}

// GetUsers returns many users by ids efficiently
func GetUsers(ctx context.Context, ids []int) ([]*models.User, []error) {
	loaders := For(ctx)
	return loaders.UserLoader.LoadMany(ctx, ids)()
}
