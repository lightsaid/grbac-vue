package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Health 健康检查，成功返回 "Success"
func CheckHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "Success"})
}

// SendGoMail 发送邮件
func SendGoMail(c *gin.Context) {

}

// UploadFiles 上传文件
func UploadFiles(c *gin.Context) {
	// c.Request.ParseMultipartForm()
}
