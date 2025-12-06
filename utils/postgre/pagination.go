package utils

import (
	"clean-arch/clean-arch-copy/app/model/postgre"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// ParsePaginationParams extracts and validates pagination parameters from query string
func ParsePaginationParams(c *fiber.Ctx) model.PaginationParams {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sortBy", "created_at")
	order := c.Query("order", "desc")
	search := c.Query("search", "")

	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	return model.PaginationParams{
		Page:   page,
		Limit:  limit,
		SortBy: sortBy,
		Order:  strings.ToLower(order),
		Search: search,
	}
}

// CalculateTotalPages calculates total pages based on total records and limit
func CalculateTotalPages(total, limit int) int {
	return (total + limit - 1) / limit
}

// CreatePaginationResponse creates a standardized pagination response
func CreatePaginationResponse(data interface{}, params model.PaginationParams, total int) fiber.Map {
	totalPages := CalculateTotalPages(total, params.Limit)

	meta := model.MetaInfo{
		Page:   params.Page,
		Limit:  params.Limit,
		Total:  total,
		Pages:  totalPages,
		SortBy: params.SortBy,
		Order:  params.Order,
		Search: params.Search,
	}

	return fiber.Map{
		"message": "Data berhasil diambil",
		"success": true,
		"data":    data,
		"meta":    meta,
	}
}
