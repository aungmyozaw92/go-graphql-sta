package models

import (
	"fmt"

	"github.com/aungmyozaw92/go-graphql/utils"
	"gorm.io/gorm"
)

// new pagination combined struct embedding + generic struct
type Cursor interface {
	GetCursor() string
}

type Edge[N Cursor] struct {
	Node   *N
	Cursor string
}

type CompositeCursor interface {
	Cursor
	Identifier
}

// fetch results for pagination
func FetchPageCompositeCursor[T CompositeCursor](dbCtx *gorm.DB,
	limit int,
	after *string,
	cursorColumn string,
	cmpOperator string,
) ([]Edge[T], *PageInfo, error) {

	nodes := make([]*T, 0)

	// order
	if cmpOperator == ">" {
		dbCtx.Order(cursorColumn + ", id")
	} else if cmpOperator == "<" {
		dbCtx.Order(cursorColumn + " DESC, id DESC")
	}

	// filter
	decodedCursor, cursorId := DecodeCompositeCursor(after)
	if decodedCursor != "" {
		dbCtx.Where(
			// [1] = column, [2] = operator
			fmt.Sprintf("%[1]s %[2]s ? OR (%[1]s = ? AND id %[2]s ?)", cursorColumn, cmpOperator),
			decodedCursor, decodedCursor, cursorId)
	}

	// db query
	dbCtx.Limit(limit + 1)
	if err := dbCtx.Find(&nodes).Error; err != nil {
		return nil, nil, err
	}

	/*
		constructing edges & page info
	*/
	count := 0
	hasNextPage := false
	edges := make([]Edge[T], 0, len(nodes))
	for _, node := range nodes {
		if count == limit {
			hasNextPage = true
		}
		if count < limit {
			var edge Edge[T]
			edge.Node = node
			edge.Cursor = EncodeCompositeCursor((*node).GetCursor(), (*node).GetId())
			edges = append(edges, edge)
			count++
		}
	}

	pageInfo := PageInfo{
		StartCursor: "",
		EndCursor:   "",
		HasNextPage: utils.NewFalse(),
	}
	if count > 0 {
		pageInfo = PageInfo{
			StartCursor: edges[0].Cursor,
			EndCursor:   edges[count-1].Cursor,
			HasNextPage: &hasNextPage,
		}
	}

	return edges, &pageInfo, nil
}