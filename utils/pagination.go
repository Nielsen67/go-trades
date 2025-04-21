package utils

import "math"

const (
	DefaultPage = 1
	DefaultSize = 5
)

type Pagination struct {
	Page      int   `json:"page" form:"page" binding:"min=1"`
	Size      int   `json:"size" form:"size" binding:"min=1,max=100"`
	Total     int64 `json:"total"`
	TotalPage int64 `json:"total_page"`
}

func GetTotalPage(totalSize int64, size int) int64 {
	return int64(math.Ceil(float64(totalSize) / float64(size)))
}
