package dto

type RoleResponse struct {
	ID          uint                  `json:"id"`
	Code        string                `json:"code"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Permissions []*PermissionResponse `json:"permissions"`
}
