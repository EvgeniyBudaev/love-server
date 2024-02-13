package searching

import (
	"fmt"
	"strings"
)

type Searching struct {
	Search string `json:"search"`
}

func ApplySearch(sqlQuery string, searchKey, searchString string) string {
	if searchString != "" {
		str := strings.ToLower(strings.TrimSpace(searchString))
		sqlQuery += fmt.Sprintf(" WHERE LOWER(%s) LIKE '%%%s%%'", searchKey, str)
	}
	return sqlQuery
}
