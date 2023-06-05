package dto

type RequestQueryParams struct {
	QueryParams
	PaginationParam
	Filter map[string]interface{}
}
