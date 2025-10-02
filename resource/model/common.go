package model

type JsonResponse struct {
	RequestId string      `json:"request_id"`
	Status    int         `json:"status_code"`
	Messages  string      `json:"messages"`
	Data      interface{} `json:"data"`
}

type JsonResponseTotal struct {
	RequestId string      `json:"request_id"`
	Status    int         `json:"status_code"`
	Messages  string      `json:"messages"`
	Total     int         `json:"total"`
	Data      interface{} `json:"data"`
}

type Pagination struct {
	Limit  int    `query:"limit" json:"limit"`
	Page   int    `query:"page" json:"page"`
	Search string `query:"search" json:"search"`
	Order  string `query:"order" json:"order"`
}

type JsonResponsError struct {
	RequestId        string      `json:"request_id"`
	StatusCode       int         `json:"status_code"`
	ErrorCode        interface{} `json:"error_code"`
	ErrorMessage     interface{} `json:"error_message"`
	DeveloperMessage interface{} `json:"developer_message"`
}

type GormWhere struct {
	Where string
	Value []interface{}
}
