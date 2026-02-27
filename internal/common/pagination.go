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
	perPage := parseIntQuery(c, "per_page", defaultPerPage)

	if page < 1 {
		page = defaultPage
	}

	if perPage < 1 {
		perPage = defaultPerPage
	}

	if perPage > maxPerPage {
		perPage = maxPerPage
	}

	offset := (page - 1) * perPage

	return Pagination{
		Page:    page,
		PerPage: perPage,
		Offset:  offset,
	}
}

func (p Pagination) ToMeta(total int64) *Meta {
	return &Meta{
		Page:    p.Page,
		PerPage: p.PerPage,
		Total:   total,
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
