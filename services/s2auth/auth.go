package s2auth

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	smbplusplus_db "github.com/babbage88/smbplusplus/database/smbplusplus_pg"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
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

type LocalAuthService struct {
	DbConn *pgxpool.Pool `json:"dbConn"`
}

func NewLoginRequest(username string, password string, isHashed bool) *UserLoginRequest {
	return &UserLoginRequest{UserName: username, Password: password}
}

func (ua *LocalAuthService) VerifyUser(userid uuid.UUID) bool {
	queries := smbplusplus_db.New(ua.DbConn)
	qry, err := queries.GetUserById(context.Background(), userid)
	if err != nil {
		slog.Error("Error querying database for user", slog.String("Error", err.Error()), slog.String("UserName", fmt.Sprint(userid)))
	}
	return qry.Enabled
}

func (us *LocalAuthService) VerifyUserPermission(ueid pgtype.UUID, permissionName string) (bool, error) {
	params := smbplusplus_db.VerifyUserPermissionByIdParams{
		UserId:     ueid,
		Permission: pgtype.Text{String: permissionName, Valid: true},
	}
	queries := smbplusplus_db.New(us.DbConn)
	qry, err := queries.VerifyUserPermissionById(context.Background(), params)
	if err != nil {
		slog.Error("error verifying user permissions", slog.String("error", err.Error()))
		return false, err
	}
	return qry, err
}

func (us *LocalAuthService) VerifyUserPermissionByRole(roleId uuid.UUID, permissionName string) (bool, error) {
	params := smbplusplus_db.VerifyUserPermissionByRoleIdParams{
		RoleId:     roleId,
		Permission: pgtype.Text{String: permissionName, Valid: true},
	}

	queries := smbplusplus_db.New(us.DbConn)
	qry, err := queries.VerifyUserPermissionByRoleId(context.Background(), params)
	if err != nil {
		slog.Error("error verifying user permissions", slog.String("error", err.Error()))
		return false, err
	}
	return qry, err
}

func (t *AuthToken) RefreshAccessTokens(dbConn *pgxpool.Pool) error {
	queries := smbplusplus_db.New(dbConn)
	qry, err := queries.GetUserById(context.Background(), t.UserID)
	if err != nil {
		slog.Error("Error querying database for user", slog.String("Error", err.Error()), slog.String("UserName", fmt.Sprint(t.UserID)))
	}
	if !qry.Enabled {
		slog.Warn("User is not enabled", slog.String("UserID", fmt.Sprint(t.UserID)))
		return fmt.Errorf("user is not enabled")
	}
	userInfo := map[string]interface{}{
		"uid":   fmt.Sprint(qry.ID),
		"email": qry.Email,
	}
	jwt_algo := os.Getenv("JWT_ALGORITHM")
	signingMethod := jwt.GetSigningMethod(jwt_algo)

	claims, err := NewInfraJWTClaims(t.UserID, userInfo)
	if err != nil {
		slog.Error("error ycreating claims", slog.String("error", err.Error()))
		return err
	}
	newAccessToken, err := NewAccessToken(claims, &signingMethod)
	if err != nil {
		slog.Error("Error creating new auth token.", slog.String("Error", err.Error()))
		return err
	}

	t.Token = newAccessToken

	return nil
}
