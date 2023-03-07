package controllers

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lightsaid/grbac/helper"
	"github.com/lightsaid/grbac/initializer"
	"github.com/lightsaid/grbac/models"
)

func CreatePermission(c *gin.Context) {
	var req models.PermissionRequest
	if ok := helper.BindRequest(c, &req); !ok {
		return
	}
	p := models.Permission{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
	}
	if err := initializer.DB.Create(&p).Error; err != nil {
		e := helper.HandleMySQLError(c, err)
		helper.ToErrResponse(c, e)
		return
	}

	helper.ToResponse(c, p)
}

func GetPermission(c *gin.Context) {
	var req models.RequestUri
	if ok := helper.BindRequestUri(c, &req); !ok {
		return
	}
	var p models.Permission
	if err := initializer.DB.Where("id = ?", req.ID).First(&p).Error; err != nil {
		e := helper.HandleMySQLError(c, err)
		helper.ToErrResponse(c, e)
		return
	}
	helper.ToResponse(c, p)
}

func UpdatePermission(c *gin.Context) {
	var uri models.RequestUri
	if ok := helper.BindRequestUri(c, &uri); !ok {
		return
	}
	var req models.PermissionRequest
	if ok := helper.BindRequest(c, &req); !ok {
		return
	}
	fmt.Println(">>> ", uri.ID)
	var p models.Permission
	if err := initializer.DB.Where("id = ?", uri.ID).First(&p).Error; err != nil {
		e := helper.HandleMySQLError(c, err)
		helper.ToErrResponse(c, e)
		return
	}
	p.BaseModel.UpdatedAt = time.Now()
	p.Name = req.Name
	p.Code = req.Code
	p.Description = req.Description
	if err := initializer.DB.Save(&p).Error; err != nil {
		e := helper.HandleMySQLError(c, err)
		helper.ToErrResponse(c, e)
		return
	}

	helper.ToResponse(c, p)
}

func DeletePermission(c *gin.Context) {
	var req models.RequestUri
	if ok := helper.BindRequestUri(c, &req); !ok {
		return
	}
	p := models.Permission{ID: req.ID}
	if err := initializer.DB.Delete(&p).Error; err != nil {
		e := helper.HandleMySQLError(c, err)
		helper.ToErrResponse(c, e)
		return
	}
	helper.ToResponse(c, nil, "操作成功")
}

func ListPermissions(c *gin.Context) {
	var req models.Pagination
	var p []*models.Permission
	if ok := helper.BindRequest(c, &req); !ok {
		return
	}
	db := initializer.DB.Offset(req.Offset()).Limit(req.Size)
	if err := db.Find(&p).Error; err != nil {
		e := helper.HandleMySQLError(c, err)
		helper.ToErrResponse(c, e)
	}
	helper.ToResponse(c, p)
}
