package middlewares

import (
	"context"
	"time"

	"github.com/aungmyozaw92/go-graphql/config"
	"github.com/aungmyozaw92/go-graphql/models"
	"github.com/gin-gonic/gin"
	"github.com/graph-gophers/dataloader/v7"
	"gorm.io/gorm"
)


type ctxKey string

const (
	loadersKey = ctxKey("dataloaders")
)

// Loaders wrap your data loaders to inject via middleware
type Loaders struct {
	UserLoader 				*dataloader.Loader[int, *models.User]
	RoleLoader 				*dataloader.Loader[int, *models.Role]
	CategoryLoader 			*dataloader.Loader[int, *models.Category]
	UnitLoader 				*dataloader.Loader[int, *models.Unit]
	ProductImageLoader      *dataloader.Loader[int, []*models.Image]
}

// NewLoaders instantiates data loaders for the middleware
func NewLoaders(conn *gorm.DB) *Loaders {
	// define the data loader
	ur := &userReader{db: conn}
	rr := &roleReader{db: conn}
	cr := &categoryReader{db: conn}
	unitr := &unitReader{db: conn}
	pImager := &imageReader{db: conn, referenceType: "products"}

	return &Loaders{
		UserLoader: dataloader.NewBatchedLoader(ur.getUsers, dataloader.WithWait[int, *models.User](time.Millisecond)),
		RoleLoader: dataloader.NewBatchedLoader(rr.getRoles, dataloader.WithWait[int, *models.Role](time.Millisecond)),
		CategoryLoader: dataloader.NewBatchedLoader(cr.getCategories, dataloader.WithWait[int, *models.Category](time.Millisecond)),
		UnitLoader: dataloader.NewBatchedLoader(unitr.getUnits, dataloader.WithWait[int, *models.Unit](time.Millisecond)),
		// ProductImageLoader: dataloader.NewBatchedLoader(pImager.GetImages, dataloader.WithWait[int, *models.Image](time.Millisecond)),
		ProductImageLoader: dataloader.NewBatchedLoader(pImager.GetImages, dataloader.WithWait[int, []*models.Image](time.Millisecond)),

	}
}

// Middleware injects data loaders into the context
// func LoaderMiddleware(conn *gorm.DB, next http.Handler) http.Handler {
// 	// return a middleware that injects the loader to the request context
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		loader := NewLoaders(conn)
// 		r = r.WithContext(context.WithValue(r.Context(), loadersKey, loader))
// 		next.ServeHTTP(w, r)
// 	})
// }
func LoaderMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		loader := NewLoaders(config.GetDB())
		ctx := context.WithValue(c.Request.Context(), loadersKey, loader)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// For returns the dataloader for a given context
func For(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}


func handleError[T any](count int, err error) []*dataloader.Result[T] {
	results := make([]*dataloader.Result[T], count)
	for i := 0; i < count; i++ {
		results[i] = &dataloader.Result[T]{Error: err}
	}
	return results
}

// T must be struct
// each id has many related results
func generateLoaderArrayResults[T models.RelatedData](results []T, referenceIds []int) (loaderResults []*dataloader.Result[[]*T]) {
	resultMap := make(map[int][]*T)
	for _, result := range results {
		// creating a new variable every turn, to avoid pointing to the adddress of result
		copy := result
		resultMap[result.GetReferenceId()] = append(resultMap[result.GetReferenceId()], &copy)
	}
	for _, id := range referenceIds {
		resultArray := resultMap[id]
		loaderResults = append(loaderResults, &dataloader.Result[[]*T]{Data: resultArray})
	}
	return loaderResults
}