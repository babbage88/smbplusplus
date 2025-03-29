package s2auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	smbplusplus_db "github.com/babbage88/smbplusplus/database/smbplusplus_pg"
	"github.com/babbage88/smbplusplus/internal/type_helper"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthService interface {
	VerifyUser(userid uuid.UUID) bool
	Login(loginReq *UserLoginRequest) UserLoginResponse
	VerifyUserPermission(executionUserId uuid.UUID, permissionsName string) (bool, error)
	CreateAuthTokenOnLogin(userid uuid.UUID, roleIds uuid.UUIDs, email string) (AuthToken, error)
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

func (us *LocalAuthService) VerifyUserPermission(ueid uuid.UUID, permissionName string) (bool, error) {
	params := smbplusplus_db.VerifyUserPermissionByIdParams{
		UserId:     pgtype.UUID{Bytes: ueid, Valid: true},
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

func (a *LocalAuthService) Login(loginReq *UserLoginRequest) UserLoginResponse {
	var response UserLoginResponse
	var result LoginResult
	username := pgtype.Text{String: loginReq.UserName, Valid: true}

	queries := smbplusplus_db.New(a.DbConn)
	qry, err := queries.GetUserLogin(context.Background(), username)
	result.PasswordValid = VerifyPassword(loginReq.Password, qry.Password.String)

	if err != nil {
		slog.Error("Error querying database for user", slog.String("UserName", loginReq.UserName))
	}

	if !result.PasswordValid {
		slog.Error("Supplied password does not match the password stored in database", slog.String("User", loginReq.UserName))
		result.Success = false
		result.Error = errors.New("password does not match")
		result.UserEnabled = qry.Enabled
		response.Result = result
		return response
	}

	if !qry.Enabled {
		slog.Error("User is disabled", slog.String("User", loginReq.UserName))
		result.Success = false
		result.UserEnabled = qry.Enabled
		result.Error = errors.New("user is diabled")
		response.Result = result
		return response
	}
	slog.Info("Login was Successful")
	result.Success = true
	result.Error = nil
	result.UserNameMatches = true

	response.Result = result
	response.UserInfo.ParseUserRowFromDb(qry)

	return response
}

func (a LocalAuthService) CreateAuthTokenOnLogin(userid uuid.UUID, roleIds uuid.UUIDs, email string) (AuthToken, error) {
	var retval AuthToken
	userInfo := map[string]interface{}{
		"email": email,
	}

	tokenString, expireTime, err := a.CreateSignedAuthTokenString(fmt.Sprint(userid), roleIds, userInfo)
	if err != nil {
		slog.Error("Error creating signed JWT token", slog.String("Error", err.Error()))
		return retval, err
	}

	retval = AuthToken{
		UserID:     userid,
		Expiration: expireTime,
		Token:      tokenString,
	}

	retval.CreateRefreshToken()

	return retval, nil
}

func (ua *LocalAuthService) CreateSignedAuthTokenString(sub string, roleIds uuid.UUIDs, userInfo interface{}) (string, time.Time, error) {
	expireMinutesEnv := os.Getenv("EXPIRATION_MINUTES")
	expireMinutes, err := type_helper.ParseInt64(expireMinutesEnv)
	if err != nil {
		slog.Error("Error parsing EXPIRATION_MINUTES, defaulting to 60.", slog.String("Error", err.Error()))
		expireMinutes = 60
	}

	jwtAlgo := os.Getenv("JWT_ALGORITHM")
	if len(jwtAlgo) < 1 {
		slog.Error("No JWT_ALGORIMTH Configered set, Setting to Default.", slog.String("Default", "HS256"))
	}
	jwtKey := []byte(os.Getenv("JWT_KEY"))

	token := jwt.New(jwt.GetSigningMethod(jwtAlgo))
	exp := time.Now().Add(time.Minute * time.Duration(expireMinutes))

	token.Claims = jwt.MapClaims{
		"sub":       sub,
		"role_ids":  roleIds,
		"user_info": userInfo,
		"exp":       exp.Unix(),
		"iss":       "smbplusplus",
	}

	signedToken, err := token.SignedString(jwtKey)
	if err != nil {
		return "", exp, err
	}

	return signedToken, exp, nil
}

func (a *LocalAuthService) VerifyToken(tokenString string) error {
	jwtKey := []byte(os.Getenv("JWT_KEY"))
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

func (ua *LocalAuthService) ParseAccessToken(accessToken string) *SmbPlusPlusClaim {
	jwtKey := os.Getenv("JWT_KEY")
	parsedAccessToken, _ := jwt.ParseWithClaims(accessToken, &SmbPlusPlusClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})

	return parsedAccessToken.Claims.(*SmbPlusPlusClaim)
}

func (ua *LocalAuthService) VerifyUserRolesForPermission(roleIds uuid.UUIDs, permissionName string) (bool, error) {
	var lastError error // Store any encountered errors for logging or debugging

	for _, roleId := range roleIds {
		hasPermission, err := ua.VerifyUserPermissionByRole(roleId, permissionName)
		if err != nil {
			// Save the error but continue checking other roles
			slog.Error("Error encountered while verifying permissions", slog.String("roleId", roleId.String()), slog.String("error", err.Error()))
			lastError = err
			continue
		}
		if hasPermission {
			return true, err
		}
	}

	if lastError != nil {
		// Log the error for debugging purposes
		slog.Error("Error occurred while verifying permissions for roles", slog.String("Error", lastError.Error()))
	}
	// Return false if no roles grant the permission
	return false, lastError
}
