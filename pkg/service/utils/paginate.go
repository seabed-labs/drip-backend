package utils

import "fmt"

func DoForPaginatedBatch(pageSize, total int, processBatch func(start, end int)) error {
	if total != 0 && pageSize == 0 {
		return fmt.Errorf("pageSize must be non-zero when total is non-zero")
	}
	page := 0
	start, end := paginate(page, pageSize, total)
	for start < end {
		processBatch(start, end)
		page++
		start, end = paginate(page, pageSize, total)
	}
	return nil
}

func paginate(pageNum int, pageSize int, sliceLength int) (int, int) {
	start := pageNum * pageSize

	if start > sliceLength {
		start = sliceLength
	}

	end := start + pageSize
	if end > sliceLength {
		end = sliceLength
	}

	return start, end
}
