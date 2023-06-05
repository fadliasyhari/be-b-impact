package response

import (
	"net/http"

	"be-b-impact.com/csr/model/dto"
	"github.com/gin-gonic/gin"
)

func SendSingleResponse(c *gin.Context, description string, data interface{}) {
	c.JSON(http.StatusOK, &SingleResponse{
		Status: Status{
			Code:        http.StatusOK,
			Description: description,
		},
		Data: data,
	})
}

func SendMultiResponse(c *gin.Context, description string, data map[string]interface{}) {
	responseData := map[string]interface{}{}

	for key, value := range data {
		responseData[key] = value
	}

	c.JSON(http.StatusOK, &SingleResponse{
		Status: Status{
			Code:        http.StatusOK,
			Description: description,
		},
		Data: responseData,
	})
}

func SendPagedResponse(c *gin.Context, description string, data []interface{}, paging dto.Paging) {
	c.JSON(http.StatusOK, &PagedResponse{
		Status: Status{
			Code:        http.StatusOK,
			Description: description,
		},
		Data:   data,
		Paging: paging,
	})
}

func SendErrorResponse(c *gin.Context, code int, description string) {
	c.AbortWithStatusJSON(code, &ErrorResponse{
		Status: Status{
			Code:        code,
			Description: description,
		},
	})
}
