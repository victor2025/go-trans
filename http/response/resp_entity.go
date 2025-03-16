package response

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go-trans/utils"
	"net/http"
)

/**
  @author: victor2022
  @since: 2025/1/5
*/

const (
	success responseStatus = "y"
	fail    responseStatus = "n"
)

type Entity struct {
	Success responseStatus `json:"success"`
	Content any            `json:"content"`
	ErrMsg  string         `json:"errMsg"`
}

type responseStatus string

func NewSuccessResponse(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Entity{Success: success, Content: data})
}

func NewFailResponse(c *gin.Context, data any, errMsg string) {
	c.JSON(http.StatusInternalServerError, Entity{Success: fail, Content: data, ErrMsg: errMsg})
}

func GetStructFromResponse(entity *Entity, target any) {
	marshal, err := json.Marshal(entity.Content)
	utils.HandleError(err)
	err = json.Unmarshal(marshal, target)
	utils.HandleError(err)
}
