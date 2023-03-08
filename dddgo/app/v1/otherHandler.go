package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type OtherHandler struct {
}

// CheckHealthHandler 健康检查，成功返回 "Success"
func (h *OtherHandler) CheckHealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "Success"})
}
