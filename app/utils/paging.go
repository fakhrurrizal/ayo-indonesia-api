package utils

import (
	"ayo-indonesia-api/app/reqres"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/guregu/null"
)

func PopulatePaging(c *gin.Context, custom string) (param reqres.ReqPaging) {
	customval := c.Query(custom)
	limit, _ := strconv.Atoi(c.Query("limit"))
	if limit == 0 {
		limit = 10
	}
	offset, _ := strconv.Atoi(c.Query("offset"))
	page, _ := strconv.Atoi(c.Query("page"))
	if page == 0 && offset == 0 {
		page = 1
		offset = 0
	}
	if page >= 1 && offset == 0 {
		offset = (page - 1) * limit
	}
	draw, _ := strconv.Atoi(c.Query("draw"))
	if draw == 0 {
		draw = 1
	}
	order := c.Query("sort")
	if strings.ToLower(order) == "asc" {
		order = "ASC"
	} else {
		order = "DESC"
	}
	sort := c.Query("order")
	if sort == "" {
		sort = "created_at " + order
	} else {
		sort = sort + " " + order + ", created_at " + order
	}

	param = reqres.ReqPaging{
		Search: c.Query("search"),
		Limit:  limit,
		Offset: offset,
		Sort:   sort,
		Custom: customval,
		Page:   page,
	}
	return
}

func PopulateResPaging(param *reqres.ReqPaging, data interface{}, totalResult int64, totalFiltered int64, lastUpdated null.Time) (output reqres.ResPaging) {
	totalPages := int(totalFiltered) / param.Limit
	if int(totalFiltered)%param.Limit > 0 {
		totalPages++
	}

	currentPage := param.Offset/param.Limit + 1
	next := false
	back := false
	if currentPage < totalPages {
		next = true
	}
	if currentPage <= totalPages && currentPage > 1 {
		back = true
	}

	output = reqres.ResPaging{
		Status:          200,
		Draw:            1,
		Data:            data,
		Search:          param.Search,
		Order:           param.Order,
		Limit:           param.Limit,
		Offset:          param.Offset,
		Sort:            param.Sort,
		Next:            next,
		Back:            back,
		TotalData:       int(totalResult),
		RecordsFiltered: int(totalFiltered),
		CurrentPage:     currentPage,
		TotalPage:       totalPages,
		LastUpdated:     lastUpdated,
	}
	return
}
