package apperror

import "fmt"

// Custom error types
var (
	ErrBlogImageNotFound    = fmt.Errorf("blog image not found")
	ErrFeatureImageNotFound = fmt.Errorf("feature image not found")
	ErrFetchingImage        = fmt.Errorf("error fetching image")
	ErrPageImageNotFound    = fmt.Errorf("page image not found")
	ErrNotFound             = fmt.Errorf("item not found")
)
