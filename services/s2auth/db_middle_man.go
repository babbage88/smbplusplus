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
