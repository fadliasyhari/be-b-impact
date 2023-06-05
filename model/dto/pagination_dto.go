package dto

type QueryParams struct {
	Query string
	Order string
	Sort  string
}

func (qp *QueryParams) IsSortValid() bool {
	return qp.Sort == "ASC" || qp.Sort == "DESC"
}

type PaginationParam struct {
	Page   int
	Limit  int
	Offset int
}

type PaginationQuery struct {
	Page int
	Take int
	Skip int
}

type Paging struct {
	Page        int
	RowsPerPage int
	TotalRows   int
	TotalPages  int
}
