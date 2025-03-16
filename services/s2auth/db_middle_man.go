package s2auth

import (
	smbplusplus_db "github.com/babbage88/smbplusplus/database/smbplusplus_pg"
	"github.com/babbage88/smbplusplus/internal/pretty"
)

func PlaceHolder() {
	pretty.Print("Need to compile to create smbplusplus_pg")
}

type DbParser interface {
	ParseUserFromDb(dbuser smbplusplus_db.User)
}

func (u *UserDao) ParseUserFromDb(dbuser smbplusplus_db.UsersWithRole) {
	u.Id = dbuser.ID
	u.UserName = dbuser.Username.String
	u.Email = dbuser.Email.String
	u.Enabled = dbuser.Enabled
	u.IsDeleted = dbuser.IsDeleted
	u.CreatedAt = dbuser.CreatedAt.Time
	u.LastModified = dbuser.LastModified.Time
	u.RoleIds = dbuser.RoleIds
	u.Roles = dbuser.Roles
}

type AuthDbParser interface {
	ParseUserFromDb(dbuser smbplusplus_db.User)
	ParseUserWithRoleFromDb(dbuser smbplusplus_db.UsersWithRole)
	ParseUserRowFromDb(dbRow smbplusplus_db.GetUserLoginRow)
	ParseAuthTokenFromDb(token smbplusplus_db.AuthToken)
	ParseUserRoleFromDb(dbRow smbplusplus_db.UserRole)
	ParseAppPermissionFromDb(dbRow smbplusplus_db.AppPermission)
	ParseRolePermissionMappingFromDb(dbRow smbplusplus_db.RolePermissionMapping)
}

func (t *AuthTokenDao) ParseAuthTokenFromDb(token smbplusplus_db.AuthToken) {
	t.Id = token.ID
	t.Token = token.Token.String
	t.UserID = token.UserID
	t.CreatedAt = token.CreatedAt.Time
	t.Expiration = token.Expiration.Time
	t.LastModified = token.LastModified.Time
}

func (ur *UserRoleDao) ParseUserRoleFromDb(dbRow smbplusplus_db.UserRole) {
	ur.Id = dbRow.ID
	ur.RoleName = dbRow.RoleName
	ur.RoleDescription = dbRow.RoleDescription.String
	ur.Enabled = dbRow.Enabled
	ur.IsDeleted = dbRow.IsDeleted
	ur.CreatedAt = dbRow.CreatedAt.Time
	ur.LastModified = dbRow.LastModified.Time
}

func (ap *AppPermissionDao) ParseAppPermissionFromDb(dbRow smbplusplus_db.AppPermission) {
	ap.Id = dbRow.ID
	ap.PermissionName = dbRow.PermissionName
	ap.PermissionDescription = dbRow.PermissionDescription.String
}

func (rpm *RolePermissionMappingDao) ParseRolePermissionMappingFromDb(dbRow smbplusplus_db.RolePermissionMapping) {
	rpm.Id = dbRow.ID
	rpm.PermissionId = dbRow.PermissionID
	rpm.RoleId = dbRow.RoleID
	rpm.Enabled = dbRow.Enabled
	rpm.CreatedAt = dbRow.CreatedAt.Time
	rpm.LastModified = dbRow.LastModified.Time
}
