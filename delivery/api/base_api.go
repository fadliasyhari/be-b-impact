package api

import (
	"be-b-impact.com/csr/delivery/api/response"
	"be-b-impact.com/csr/model/dto"
	"github.com/gin-gonic/gin"
)

type BaseApi struct{}

func (b *BaseApi) ParseRequestBody(c *gin.Context, payload interface{}) error {
	if err := c.ShouldBindJSON(payload); err != nil {
		return err
	}
	return nil
}

func (b *BaseApi) NewSuccessSingleResponse(c *gin.Context, description string, data interface{}) {
	response.SendSingleResponse(c, description, data)
}

func (b *BaseApi) NewSuccessMultiResponse(c *gin.Context, description string, data map[string]interface{}) {
	response.SendMultiResponse(c, description, data)
}

func (b *BaseApi) NewSuccessPagedResponse(c *gin.Context, description string, data []interface{}, paging dto.Paging) {
	response.SendPagedResponse(c, description, data, paging)
}

func (b *BaseApi) NewFailedResponse(c *gin.Context, code int, description string) {
	response.SendErrorResponse(c, code, description)
}
