package utils

import (
	"ayo-indonesia-api/app/reqres"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

func PopulatePaging(c *gin.Context, custom string) reqres.ReqPaging {
	customval := c.Query(custom)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	sort := c.DefaultQuery("sort", "id")
	order := c.DefaultQuery("order", "desc")
	search := c.Query("search")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	if order != "asc" && order != "desc" {
		order = "desc"
	}

	offset := (page - 1) * limit

	return reqres.ReqPaging{
		Page:   page,
		Limit:  limit,
		Sort:   sort,
		Order:  order,
		Search: search,
		Offset: offset,
		Custom: customval,
	}
}

func PopulateResPaging(param *reqres.ReqPaging, data interface{}, total int64, totalFiltered int64) reqres.ResPaging {
	totalPages := int(math.Ceil(float64(totalFiltered) / float64(param.Limit)))

	return reqres.ResPaging{
		Status:        200,
		Data:          data,
		TotalData:     total,
		TotalFiltered: totalFiltered,
		Page:          param.Page,
		Limit:         param.Limit,
		TotalPages:    totalPages,
	}
}
