package directives

import (
	"context"
	"errors"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/aungmyozaw92/go-graphql/config"
	"github.com/aungmyozaw92/go-graphql/middlewares"
	"github.com/aungmyozaw92/go-graphql/models"
	"github.com/aungmyozaw92/go-graphql/utils"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func Auth(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	
	tokenData := middlewares.CtxValue(ctx)
	// fmt.Println("tokenData", tokenData)
	if tokenData == nil {
		return nil, &gqlerror.Error{
			Message: "Access Denied",
		}
	}

	userId, ok := utils.GetUserIdFromContext(ctx)
	if !ok || userId == 0 {
		return nil, &gqlerror.Error{
			Message: "Access Denied",
		}
	}

	user, err := models.GetUser(ctx, userId)
	if err != nil {
		return nil, &gqlerror.Error{
			Message: err.Error(),
		}
	}
	if !*user.IsActive {
		return nil, &gqlerror.Error{
			Message: "User is disabled",
		}
	}

	gqlpath := graphql.GetPath(ctx).String()
	
	// user is either owner or custom
	if err := authorizeUser(ctx, user.RoleId, gqlpath); err != nil {
		return nil, &gqlerror.Error{
			Message: err.Error(),
		}
	}
	
	ctx = context.WithValue(ctx, utils.ContextKeyUsername, user.Username)

	return next(ctx)
}

// retrieve role's allowed query paths from redis and check if the gqlpath is allowed
func authorizeUser(ctx context.Context, roleId int, gqlpath string) error {
	
	var queryPaths map[string]bool
	exists, err := config.GetRedisObject("AllowedPaths:Role:"+fmt.Sprint(roleId), &queryPaths)
	if err != nil {
		return err
	}

	if !exists {

		queryPaths, err = models.GetQueryPathsFromRole(ctx, roleId)
		if err != nil {
			return err
		}

		// store in redis
		if err := config.SetRedisObject("AllowedPaths:Role:"+fmt.Sprint(roleId), &queryPaths, 0); err != nil {
			return err
		}
	}

	// check if current path is allowed for current user
	// using a map for faster look up, non-existent key will return false, default zero for boolean
	fmt.Println(queryPaths)
	fmt.Println(gqlpath)
	if allowed := queryPaths[gqlpath]; !allowed {
	if gqlpath == "logout" {
		return nil
	}
		return errors.New("permission not allow")
	}

	return nil
}
