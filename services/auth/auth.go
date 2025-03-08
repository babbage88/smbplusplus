package s2auth

import "time"

type AppPermissionDao struct {
	Id                    int32  `json:"id"`
	PermissionName        string `json:"permissionName"`
	PermissionDescription string `json:"permissionDescription"`
}

type RolePermissionMappingDao struct {
	Id           int32     `json:"id"`
	RoleId       int32     `json:"roleId"`
	PermissionId int32     `json:"permissionId"`
	Enabled      bool      `json:"enabled"`
	CreatedAt    time.Time `json:"createdAt"`
	LastModified time.Time `json:"lastModified"`
}

type UserRoleDao struct {
	Id              int32     `json:"id"`
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
	Id           int32     `json:"id"`
	UserName     string    `json:"username"`
	Email        string    `json:"email"`
	Roles        []string  `json:"roles"`
	RoleIds      []int32   `json:"role_ids"`
	CreatedAt    time.Time `json:"createdAt"`
	LastModified time.Time `json:"lastModified"`
	Enabled      bool      `json:"enabled"`
	IsDeleted    bool      `json:"isDeleted"`
}

type AuthTokenDao struct {
	Id           int32     `json:"id"`
	UserID       int32     `json:"user_id"`
	Token        string    `json:"token"`
	Expiration   time.Time `json:"expiration"`
	CreatedAt    time.Time `json:"created_at"`
	LastModified time.Time `json:"last_modified"`
}
