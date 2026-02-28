package dtos

type PaginationRequest struct {
	Limit  int `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset int `form:"offset" binding:"omitempty,min=0"`
}

type DateRangeRequest struct {
	StartDate *string `form:"start_date"`
	EndDate   *string `form:"end_date"`
}

type PaginatedResponse struct {
	Total int         `json:"total"`
	Data  interface{} `json:"data"`
}
