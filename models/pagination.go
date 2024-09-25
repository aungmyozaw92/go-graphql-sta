package models

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
)

type PageInfo struct {
	StartCursor string `json:"startCursor"`
	EndCursor   string `json:"endCursor"`
	HasNextPage *bool  `json:"hasNextPage,omitempty"`
}

func DecodeCompositeCursor(cursor *string) (string, int) {
	if cursor == nil || *cursor == "" {
		return "", 0
	}

	decoded, err := base64.StdEncoding.DecodeString(*cursor)
	if err != nil {
		return "", 0
	}

	parts := strings.Split(string(decoded), "|")
	if len(parts) != 2 {
		return "", 0
	}

	id, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", 0
	}

	return parts[0], id
}

func EncodeCompositeCursor(transactionDateTime string, id int) string {
	cursor := fmt.Sprintf("%s|%d", transactionDateTime, id)
	return base64.StdEncoding.EncodeToString([]byte(cursor))
}