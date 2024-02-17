package pagination

import (
	"context"
	"database/sql"
	"fmt"
)

type Pagination struct {
	HasNext     bool   `json:"hasNext"`
	HasPrevious bool   `json:"hasPrevious"`
	CountPages  uint64 `json:"countPages"`
	Size        uint64 `json:"size"`
	Page        uint64 `json:"page"`
	TotalItems  uint64 `json:"totalItems"`
}

func NewPagination(pag *Pagination) *Pagination {
	return &Pagination{
		HasNext:     pag.HasNext,
		HasPrevious: pag.HasPrevious,
		CountPages:  pag.CountPages,
		Size:        pag.Size,
		Page:        pag.Page,
		TotalItems:  pag.TotalItems,
	}
}

func ApplyPagination(sqlQuery string, page uint64, size uint64) string {
	offset := (page - 1) * size
	sqlQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", size, offset)
	return sqlQuery
}

func GetTotalItems(ctx context.Context, db *sql.DB, sqlQuery string, args ...interface{}) (uint64, error) {
	var totalItems uint64
	err := db.QueryRowContext(ctx, sqlQuery, args...).Scan(&totalItems)
	if err != nil {
		return 0, err
	}
	return totalItems, nil
}

func getCountPages(size uint64, totalItems uint64) uint64 {
	return (totalItems + size - 1) / size
}

func GetPagination(size uint64, page uint64, totalItems uint64) *Pagination {
	hasPrevious := page > 1
	hasNext := (page * size) < totalItems
	constPages := getCountPages(size, totalItems)
	paging := NewPagination(&Pagination{
		HasNext:     hasNext,
		HasPrevious: hasPrevious,
		CountPages:  constPages,
		Size:        size,
		Page:        page,
		TotalItems:  totalItems,
	})
	return paging
}
