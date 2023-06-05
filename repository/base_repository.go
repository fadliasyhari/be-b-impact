package repository

import (
	"fmt"

	"be-b-impact.com/csr/model/dto"
	"be-b-impact.com/csr/utils/common"
)

type BaseRepository[T any] interface {
	Search(by map[string]interface{}) ([]T, error)
	List() ([]T, error)
	Get(id string) (*T, error)
	Save(payload *T) error
	Update(payload *T) error
	Delete(id string) error
}

type BaseRepositoryEmailPhone[T any] interface {
	GetByEmail(email string) (*T, error)
	GetByPhone(phone string) (*T, error)
}

type BaseRepositoryCount[T any] interface {
	CountData(fieldname string, id string) error
}
type BaseRepositoryPaging[T any] interface {
	Paging(requestQueryParam dto.RequestQueryParams) ([]T, dto.Paging, error)
}

func pagingValidate(requestQueryParam dto.RequestQueryParams) (dto.PaginationQuery, string) {
	var paginationQuery = common.GetPaginationParams(requestQueryParam.PaginationParam)
	orderQuery := "id"
	if requestQueryParam.QueryParams.Order != "" && requestQueryParam.QueryParams.Sort != "" {
		sorting := "ASC"
		if requestQueryParam.QueryParams.Sort == "DESC" {
			sorting = "DESC"
		}
		orderQuery = fmt.Sprintf("%s %s", requestQueryParam.QueryParams.Order, sorting)
	}
	return paginationQuery, orderQuery
}
