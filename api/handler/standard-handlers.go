package handler

import (
	"Study/websocket/models"
	"fmt"

	"github.com/gin-gonic/gin"
)


func (h *handler) handleErrorResponse(c *gin.Context, code int, message string, err interface{}) {
	c.JSON(code, models.ErrorModel{
		Code:    fmt.Sprint(code),
		Message: message,
		Error:   err,
	})
}

