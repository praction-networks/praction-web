// PaginationParams represents pagination, filtering, and sorting parameters
package utils

type PaginationParams struct {
	Page      int
	PageSize  int
	SortField string
	SortOrder int
	Filters   map[string]interface{}
}
