package common

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Pagination struct {
	Page    int
	PerPage int
	Offset  int
}

const (
	defaultPage    = 1
	defaultPerPage = 20
	maxPerPage     = 100
)

func NewPagination(c *gin.Context) Pagination {
	page := parseIntQuery(c, "page", defaultPage)
	limit := parseIntQuery(c, "limit", 0)
	if limit <= 0 {
		limit = parseIntQuery(c, "per_page", defaultPerPage)
	}

	if page < 1 {
		page = defaultPage
	}

	if limit < 1 {
		limit = defaultPerPage
	}

	if limit > maxPerPage {
		limit = maxPerPage
	}

	offset := (page - 1) * limit

	return Pagination{
		Page:    page,
		PerPage: limit,
		Offset:  offset,
	}
}

func (p Pagination) ToMeta(total int64) *Meta {
	return &Meta{
		Page:        p.Page,
		PerPage:     p.PerPage,
		Total:       total,
		HasPrevious: p.Page > 1,
		HasNext:     total > int64(p.Page*p.PerPage),
	}
}

func parseIntQuery(c *gin.Context, key string, defaultVal int) int {
	val := c.Query(key)
	if val == "" {
		return defaultVal
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return i
}
