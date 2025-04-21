package utils

import (
	"math"
	"strconv"
)

const DEFAULT_PAGE_SIZE = 10

func ExtractPaginationMetadata(pageParam, pageSizeParam string) (int, int) {
	pageSize := DEFAULT_PAGE_SIZE
	page := 0

	if parsedPage, parsedPageErr := strconv.Atoi(pageParam); parsedPageErr == nil {
		page = parsedPage
	}
	if parsedPageSize, parsedPageSizeErr := strconv.Atoi(pageSizeParam); parsedPageSizeErr == nil {
		pageSize = parsedPageSize
	}
	return page, pageSize
}

func CalculatePaginationTotalPage(totalCount int, pageSize int) int {
	return int(math.Ceil(float64(totalCount) / float64(pageSize)))
}
