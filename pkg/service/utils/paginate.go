package utils

import "fmt"

func DoForPaginatedBatch(pageSize, total int, processBatch func(start, end int) error, handleError func(err error) error) error {
	if total != 0 && pageSize == 0 {
		return fmt.Errorf("pageSize must be non-zero when total is non-zero")
	}
	page := 0
	start, end := Paginate(page, pageSize, total)
	for start < end {
		if err := processBatch(start, end); err != nil {
			if err := handleError(err); err != nil {
				return err
			}
		}
		page++
		start, end = Paginate(page, pageSize, total)
	}
	return nil
}

func Paginate(pageNum int, pageSize int, sliceLength int) (int, int) {
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
