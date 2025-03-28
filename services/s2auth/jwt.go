package s2auth

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/babbage88/smbplusplus/internal/type_helper"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func NewInfraJWTClaims(userid uuid.UUID, userInfo interface{}) (SmbPlusPlusClaim, error) {
	env_expire_minutes := os.Getenv("EXPIRATION_MINUTES")
	expire_minutes, err := type_helper.ParseInt64(env_expire_minutes)
	if err != nil {
		slog.Error("Error Parsing int64 from .env EXPIRATION_MINUTES, setting value to 60.", slog.String("Error", err.Error()))
		expire_minutes = 60
	}

	exp := time.Now().Add(time.Minute * time.Duration(expire_minutes))
	retVal := &SmbPlusPlusClaim{
		&jwt.RegisteredClaims{
			// Set the userid and expiration as the standard claim.
			Issuer:    "goinfra",
			ExpiresAt: jwt.NewNumericDate(exp),
			Subject:   fmt.Sprint(userid),
		},
		// UserInfo passed from caller as map[string]string
		userInfo,
	}
	return *retVal, nil
}

func NewAccessToken(claims SmbPlusPlusClaim, algo *jwt.SigningMethod) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return accessToken.SignedString([]byte(os.Getenv("JWT_KEY")))
}

func NewRefreshToken(claims jwt.RegisteredClaims) (string, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return refreshToken.SignedString([]byte(os.Getenv("JWT_KEY")))
}

func ParseAccessToken(accessToken string) *SmbPlusPlusClaim {
	parsedAccessToken, _ := jwt.ParseWithClaims(accessToken, &SmbPlusPlusClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_KEY")), nil
	})

	return parsedAccessToken.Claims.(*SmbPlusPlusClaim)
}

func ParseRefreshToken(refreshToken string) *jwt.RegisteredClaims {
	parsedRefreshToken, _ := jwt.ParseWithClaims(refreshToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_KEY")), nil
	})

	return parsedRefreshToken.Claims.(*jwt.RegisteredClaims)
}

func (t *AuthToken) CreateRefreshToken() {
	jwtKey := os.Getenv("JWT_KEY")
	//refreshEnvValue, err := type_helper.ParseInt64(os.Getenv("REFRESH_EXPIRATION_HOURS"))
	refreshExpiration := time.Now().Add(time.Hour * 48).Unix()
	refreshToken := jwt.New(jwt.SigningMethodHS256)

	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["sub"] = t.UserID
	rtClaims["exp"] = refreshExpiration

	rt, err := refreshToken.SignedString([]byte(jwtKey))
	if err != nil {
		slog.Error("Error signing refresh token", slog.String("Error", err.Error()))
	}

	t.RefreshToken = rt
}
