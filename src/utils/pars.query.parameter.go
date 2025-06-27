package utils

import (
	"net/url"
	"strconv"
)

func ParseQueryParams(query url.Values) (PaginationParams, error) {
	page, err := strconv.Atoi(query.Get("Page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(query.Get("PageSize"))
	if err != nil || pageSize < 1 {
		pageSize = 500
	}

	sortField := query.Get("sortField")
	if sortField == "" {
		sortField = "createdAt"
	}

	sortOrder := -1 // Default to DESC
	if sortOrderStr := query.Get("sortOrder"); sortOrderStr != "" {
		if parsedOrder, err := strconv.Atoi(sortOrderStr); err == nil {
			if parsedOrder == 1 || parsedOrder == -1 {
				sortOrder = parsedOrder
			}
		}
	}

	// Parse other filters
	filters := make(map[string]interface{})
	for key, values := range query {
		if key == "Page" || key == "PageSize" || key == "sortField" || key == "sortOrder" {
			continue
		}
		if len(values) > 1 {
			filters[key] = values
		} else {
			filters[key] = values[0]
		}
	}

	return PaginationParams{
		Page:      page,
		PageSize:  pageSize,
		SortField: sortField,
		SortOrder: sortOrder,
		Filters:   filters,
	}, nil
}
