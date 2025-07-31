package reqres

type ReqPaging struct {
	Page   int         `json:"page"`
	Limit  int         `json:"limit"`
	Sort   string      `json:"sort"`
	Order  string      `json:"order"`
	Search string      `json:"search"`
	Offset int         `json:"offset"`
	Custom interface{} `default:""`
}

type ResPaging struct {
	Data          interface{} `json:"data"`
	TotalData     int64       `json:"total_data"`
	TotalFiltered int64       `json:"total_filtered"`
	Page          int         `json:"page"`
	Limit         int         `json:"limit"`
	TotalPages    int         `json:"total_pages"`
	Status        int         `json:"status"`
}
