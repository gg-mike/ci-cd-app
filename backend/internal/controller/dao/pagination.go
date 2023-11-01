package dao

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Paginate(ctx *gin.Context) (int, int, string, error) {
	var page, size int
	var err error
	page_, ok := ctx.GetQuery("page")
	if !ok {
		page = -1
	} else {
		page, err = strconv.Atoi(page_)
		if err != nil {
			return -1, -1, "", errors.New(`error parsing query param "page"`)
		}
	}

	size_, ok := ctx.GetQuery("size")
	if !ok {
		size = -1
	} else {
		size, err = strconv.Atoi(size_)
		if err != nil {
			return -1, -1, "", errors.New(`error parsing query param "size"`)
		}
	}

	order, _ := ctx.GetQuery("order")

	if page >= 0 && size < 0 {
		return -1, -1, "", errors.New("page size cannot be lower than 0 for nonnegative page number")
	}
	if size == 0 {
		return -1, -1, "", errors.New("page size cannot be 0")
	}
	if page < 0 && size < 0 {
		return -1, -1, order, nil
	}
	if page < 0 && size >= 0 {
		return -1, size, order, nil
	}

	return page * size, size, order, nil
}
