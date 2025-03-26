package s2auth

import (
	"time"

	"github.com/google/uuid"
)

type AuthService interface {
	VerifyUser(userid uuid.UUID) bool
	RefreshAuthTokens() (AuthTokenDao, error)
	VerifyUserPermission(executionUserId uuid.UUID, permissionsName string) (bool, error)
	NewLoginRequest(username string, password string, isHashed bool) *UserLoginResponse
	CreateAuthToken(userid uuid.UUID, role string, email string) (AuthTokenDao, error)
	CreateSignedTokenString(sub string, userInfo interface{}) (string, time.Time, error)
	VerifyToken(tokenString string) error
	VerifyUserRolesForPermission(roleIds uuid.UUIDs, permissionName string) (bool, error)
	VerifyUserPermissionByRole(roleId uuid.UUID, permissionName string) (bool, error)
}
