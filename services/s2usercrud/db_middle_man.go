package s2usercrud

import (
	smbplusplus_db "github.com/babbage88/smbplusplus/database/smbplusplus_pg"
)

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

func (u *UserDao) ParseUserWithRolesFromDb(dbuser smbplusplus_db.UsersWithRole) {
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

func (u *UserDao) ParseUserFromDb(dbuser smbplusplus_db.User) {
	u.Id = dbuser.ID
	u.UserName = dbuser.Username.String
	u.Email = dbuser.Email.String
	u.Enabled = dbuser.Enabled
	u.IsDeleted = dbuser.IsDeleted
	u.CreatedAt = dbuser.CreatedAt.Time
	u.LastModified = dbuser.LastModified.Time
}

func (u *UserDao) ParseUserRowFromDb(dbuser smbplusplus_db.GetUserLoginRow) {
	u.Id = dbuser.ID
	u.UserName = dbuser.Username.String
	u.Email = dbuser.Email.String
	u.Enabled = dbuser.Enabled
	u.RoleIds = dbuser.RoleIds
	u.Roles = dbuser.Roles
}
