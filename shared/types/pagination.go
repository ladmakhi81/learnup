package types

type PaginationRes struct {
	Row         any `json:"row"`
	CurrentPage int `json:"currentPage"`
	TotalPage   int `json:"totalPage"`
	TotalCount  int `json:"totalCount"`
}

func NewPaginationRes(row any, currentPage, totalPage, totalCount int) PaginationRes {
	return PaginationRes{
		Row:         row,
		CurrentPage: currentPage + 1,
		TotalPage:   totalPage,
		TotalCount:  totalCount,
	}
}
