package pgrepos

import (
	"context"
	"fmt"
	"gorm.io/gorm"
)

func GetQueryExecutor(ctx context.Context) (*gorm.DB, error) {
	queryExecutor := ctx.Value("QueryExecutor")
	if queryExecutor == nil {
		return nil, fmt.Errorf("QueryExecutor was not set")
	}
	db, ok := queryExecutor.(*gorm.DB)
	if !ok {
		return nil, fmt.Errorf("QueryExecutor is of invalid type")
	}
	return db, nil
}