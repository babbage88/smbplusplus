package s2usercrud

import (
	"time"

	"github.com/google/uuid"
)

type AppPermissionDao struct {
	Id                    uuid.UUID `json:"id"`
	PermissionName        string    `json:"permissionName"`
	PermissionDescription string    `json:"permissionDescription"`
}

type RolePermissionMappingDao struct {
	Id           uuid.UUID `json:"id"`
	RoleId       uuid.UUID `json:"roleId"`
	PermissionId uuid.UUID `json:"permissionId"`
	Enabled      bool      `json:"enabled"`
	CreatedAt    time.Time `json:"createdAt"`
	LastModified time.Time `json:"lastModified"`
}

type UserRoleDao struct {
	Id              uuid.UUID `json:"id"`
	RoleName        string    `json:"roleName"`
	RoleDescription string    `json:"roleDesc"`
	Enabled         bool      `json:"enabled"`
	IsDeleted       bool      `json:"isDeleted"`
	CreatedAt       time.Time `json:"createdAt"`
	LastModified    time.Time `json:"lastModified"`
}

// Respose will return login result and the user info.
// swagger:response UserDao
// This text will appear as description of your response body.
type UserDao struct {
	// in:body
	Id           uuid.UUID  `json:"id"`
	UserName     string     `json:"username"`
	Email        string     `json:"email"`
	Roles        []string   `json:"roles"`
	RoleIds      uuid.UUIDs `json:"role_ids"`
	CreatedAt    time.Time  `json:"createdAt"`
	LastModified time.Time  `json:"lastModified"`
	Enabled      bool       `json:"enabled"`
	IsDeleted    bool       `json:"isDeleted"`
}
