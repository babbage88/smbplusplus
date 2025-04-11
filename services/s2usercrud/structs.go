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

// Create User Request takes  in Username, Password, Email, and Role fo the new user.
// swagger:parameters idOfcreateUserEndpoint
type CreateNewUserReqWrapper struct {
	// in:body
	Body CreateNewUserRequest `json:"body"`
}

type CreateNewUserRequest struct {
	NewUsername     string `json:"newUsername"`
	NewUserPassword string `json:"newPassword"`
	NewUserEmail    string `json:"newEmail"`
}

type GetAllUsersResponse struct {
	Users []UserDao `json:"users"`
}

// swagger:parameters idOfgetUserByIdEndpoint
type GetUserByIdRequest struct {
	// ID of user
	//
	// In: path
	ID string `json:"ID"`
}

type GetUserByIdResponse struct {
	User UserDao `json:"user"`
}

type GetAllRolesResponse struct {
	UserRoles []UserRoleDao `json:"userRoles"`
}

type GetAllAppPermissionsResponse struct {
	AppPermissions []AppPermissionDao `json:"appPermissions"`
}

// Allows and Admin user to update another user's password.
// swagger:parameters idOfUpdateUserPw
type UpdatePasswordRequestWrapper struct {
	// in:body
	Body UpdateUserPasswordRequest `json:"body"`
}

type UpdateUserPasswordRequest struct {
	TargetUserId uuid.UUID `json:"targetUserId"`
	NewPassword  string    `json:"newPassword"`
}

type UserPasswordUpdateResponse struct {
	Success      bool      `json:"success"`
	Error        error     `json:"error"`
	TargetUserId uuid.UUID `json:"targetUserId"`
}

// User Id of the executing and target users
// swagger:parameters idOfEnableUser
type EnableUserRequestWrapper struct {
	//in:body
	Body EnableUserRequest `json:"body"`
}

// User Id of the executing and target users
// swagger:parameters idOfDisableUser
type DisableUserRequestWrapper struct {
	//in:body
	Body DisableUserRequest `json:"body"`
}

type DisableUserRequest struct {
	TargetUserId uuid.UUID `json:"targetUserId"`
}

type EnableUserRequest struct {
	TargetUserId uuid.UUID `json:"targetUserId"`
}

type EnableDisableUserResponse struct {
	ModifiedUserInfo UserDao `json:"modifiedUserInfo"`
	Error            error   `json:"error"`
}

// User Id of the executing and target users
// swagger:parameters idOfUpdateUserRole idOfdisableUserRoleMapping
type UpdateUserRoleMappingRequestWrapper struct {
	//in:body
	Body UpdateUserRoleMappingRequest `json:"body"`
}

type UpdateUserRoleMappingRequest struct {
	TargetUserId uuid.UUID `json:"targetUserId"`
	RoleId       uuid.UUID `json:"roleId"`
}

type UpdateUserRoleMappingResponse struct {
	Success bool  `json:"success"`
	Error   error `json:"error"`
}

// User Id of the executing and target users
// swagger:parameters idOfCreateUserRole
type CreateUserRoleRequestWrapper struct {
	//in:body
	Body CreateUserRoleRequest `json:"body"`
}

type CreateUserRoleRequest struct {
	RoleName        string `json:"roleName"`
	RoleDescription string `json:"roleDesc"`
}

type CreateUserRoleResponse struct {
	Error           error       `json:"error"`
	NewUserRoleInfo UserRoleDao `json:"newUserRoleInfo"`
}

// Name and Description for new App Permission
// swagger:parameters idOfCreateAppPermission
type CreateAppPermissionRequestWrapper struct {
	//in: body
	Body CreateAppPermissionRequest `json:"body"`
}
type CreateAppPermissionRequest struct {
	PermissionName        string `json:"name"`
	PermissionDescription string `json:"descripiton"`
}

type CreateAppPermissionResponse struct {
	NewAppPermissionInfo AppPermissionDao `json:"newPermissionInfo"`
	Error                error            `json:"error"`
}

// Name and Description for new App Permission
// swagger:parameters idOfCreateRolePermissionMapping
type CreateRolePermissionMappingRequestWrapper struct {
	//in: body
	Body CreateRolePermissionMappingRequest `json:"body"`
}

type CreateRolePermissionMappingRequest struct {
	RoleId       uuid.UUID `json:"roleId"`
	PermissionId uuid.UUID `json:"permId"`
}

type CreateRolePermissionMappingResponse struct {
	NewMappingInfo RolePermissionMappingDao `json:"newMappingInfo"`
	Error          error                    `json:"error"`
}

// Mark user as deleted in Database. Will no longer show in UI unless explicityly restored
// swagger:parameters idOfSoftDeleteUserById
type SoftDeleteUserByIdRequestWrapper struct {
	//in: body
	Body SoftDeleteUserByIdRequest `json:"body"`
}

type SoftDeleteUserByIdRequest struct {
	TargetUserId uuid.UUID `json:"targetUserId"`
}

type SoftDeleteUserByIdResponse struct {
	DeletedUserInfo UserDao `json:"deletedUserInfo"`
	Error           error   `json:"error"`
}
