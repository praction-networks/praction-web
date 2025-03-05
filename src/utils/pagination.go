// PaginationParams represents pagination, filtering, and sorting parameters
package utils

type PaginationParams struct {
	Page      int                    // Current page number
	PageSize  int                    // Number of items per page
	SortField string                 // Field to sort by
	SortOrder int                    // 1 for ascending, -1 for descending
	Filters   map[string]interface{} // Filters for query
}
