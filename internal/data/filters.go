package data

import (
	"strings"

	"github.com/blessedmadukoma/greenlight/internal/validator"
)

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafelist []string
}

func ValidateFilters(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 10_000_000, "page", "must be a maximum of 10,000,000")
	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")

	v.Check(validator.In(f.Sort, f.SortSafelist...), "sort", "invalid sort value")
}

// sortColumn returns the SQL column name for the provided sort key.
func (f Filters) sortColumn() string {
	for _, safeValue := range f.SortSafelist {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}

	panic("unsafe sort parameter: " + f.Sort)
}

// sortDirection returns the SQL sort direction ("ASC" or "DESC") for the provided sort key.
func (f Filters) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}

	return "ASC"
}

// limit	returns the LIMIT value for a SQL query based on the page size.
func (f Filters) limit() int {
	return f.PageSize
}

// offset returns the OFFSET value for a SQL query based on the page number and page size.
func (f Filters) offset() int {
	return (f.Page - 1) * f.PageSize
}
