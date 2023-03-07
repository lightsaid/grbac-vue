package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/lightsaid/grbac/helper"
	"github.com/lightsaid/grbac/initializer"
	"github.com/lightsaid/grbac/models"
)

func CreateRole(c *gin.Context) {
	var req models.RoleRequest
	if ok := helper.BindRequest(c, &req); !ok {
		return
	}
	role := models.Role{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
	}

	permissions := make([]*models.Permission, len(req.PermissionIds))
	for i, _ := range req.PermissionIds {
		permissions[i] = &models.Permission{ID: req.PermissionIds[i]}
	}
	role.Permissions = permissions

	if err := initializer.DB.Create(&role).Error; err != nil {
		e := helper.HandleMySQLError(c, err)
		helper.ToErrResponse(c, e)
		return
	}

	helper.ToResponse(c, role)
}

func GetRole(c *gin.Context) {
	var req models.RequestUri
	if ok := helper.BindRequestUri(c, &req); !ok {
		return
	}
	var role models.Role

	if err := initializer.DB.Preload("Permissions").Where("id = ?", req.ID).First(&role).Error; err != nil {
		e := helper.HandleMySQLError(c, err)
		helper.ToErrResponse(c, e)
		return
	}

	helper.ToResponse(c, role)
}

func UpdateRole(c *gin.Context) {
	var uri models.RequestUri
	if ok := helper.BindRequestUri(c, &uri); !ok {
		return
	}
	var req models.RoleRequest
	if ok := helper.BindRequest(c, &req); !ok {
		return
	}

	// 先删除关系表
	var result interface{}
	if err := initializer.DB.Table("role_permissions").Where("role_id", uri.ID).Delete(&result).Error; err != nil {
		e := helper.HandleMySQLError(c, err)
		helper.ToErrResponse(c, e)
		return
	}

	role := models.Role{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
	}

	permissions := make([]*models.Permission, len(req.PermissionIds))
	for i, _ := range req.PermissionIds {
		permissions[i] = &models.Permission{ID: req.PermissionIds[i]}
	}
	role.Permissions = permissions
	// 更新
	if err := initializer.DB.Model(&role).Updates(role).Error; err != nil {
		e := helper.HandleMySQLError(c, err)
		helper.ToErrResponse(c, e)
		return
	}

	helper.ToResponse(c, role)
}

func DeleteRole(c *gin.Context) {
	var req models.RequestUri
	if ok := helper.BindRequestUri(c, &req); !ok {
		return
	}
	role := models.Role{ID: req.ID}
	if err := initializer.DB.Delete(&role).Error; err != nil {
		e := helper.HandleMySQLError(c, err)
		helper.ToErrResponse(c, e)
		return
	}
	helper.ToResponse(c, nil, "操作成功")
}

func ListRoles(c *gin.Context) {
	var req models.Pagination
	var role []*models.Role
	if ok := helper.BindRequest(c, &req); !ok {
		return
	}
	db := initializer.DB.Offset(req.Offset()).Limit(req.Size)
	if err := db.Preload("Permissions").Find(&role).Error; err != nil {
		e := helper.HandleMySQLError(c, err)
		helper.ToErrResponse(c, e)
	}
	helper.ToResponse(c, role)
}
