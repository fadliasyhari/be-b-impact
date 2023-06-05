package response

import "be-b-impact.com/csr/model/dto"

type Status struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
}

type ErrorResponse struct {
	Status Status `json:"status"`
}

type SingleResponse struct {
	Status Status      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
}

type PagedResponse struct {
	Status Status        `json:"status"`
	Data   []interface{} `json:"data,omitempty"`
	Paging dto.Paging    `json:"paging,omitempty"`
}
