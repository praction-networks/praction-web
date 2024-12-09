package utils

import (
	"net/url"
	"strconv"
)

// ParseQueryParams parses the query parameters from the URL and returns pagination, sorting, and filters
func ParseQueryParams(query url.Values) (PaginationParams, error) {
	page, err := strconv.Atoi(query.Get("Page"))
	if err != nil || page < 1 {
		page = 1 // Default to page 1
	}

	pageSize, err := strconv.Atoi(query.Get("PageSize"))
	if err != nil || pageSize < 1 {
		pageSize = 10 // Default to 10 items per page
	}

	sortField := query.Get("sortField")
	if sortField == "" {
		sortField = "createdAt" // Default sorting field
	}

	sortOrder := 1 // Default to ascending
	if query.Get("sortOrder") == "desc" {
		sortOrder = -1
	}

	// Parse filters dynamically
	filters := make(map[string]string)
	for key, values := range query {
		// Skip reserved parameters (pagination and sorting)
		if key == "Page" || key == "PageSize" || key == "sortField" || key == "sortOrder" {
			continue
		}
		// Use the first value for simplicity
		if len(values) > 0 {
			filters[key] = values[0]
		}
	}

	// Return the parsed parameters
	return PaginationParams{
		Page:      page,
		PageSize:  pageSize,
		SortField: sortField,
		SortOrder: sortOrder,
		Filters:   filters,
	}, nil
}
