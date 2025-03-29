package s2auth

import (
	"time"

	"github.com/babbage88/smbplusplus/services/s2usercrud"
	"github.com/golang-jwt/jwt/v5"

	"github.com/google/uuid"
)

// Respose will return login result and the user info.
// swagger:response AuthToken
// This text will appear as description of your response body.
type AuthToken struct {
	// in:body
	UserID       uuid.UUID `json:"user_id"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refreshToken"`
	Expiration   time.Time `json:"expiration"`
}

type AuthTokenDao struct {
	Id           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	Token        string    `json:"token"`
	Expiration   time.Time `json:"expiration"`
	CreatedAt    time.Time `json:"created_at"`
	LastModified time.Time `json:"last_modified"`
}

// Login Request takes  in Username and Password.
// swagger:parameters idOfloginEndpoint
type UserLoginReqWrapper struct {
	// in:body
	Body UserLoginRequest `json:"body"`
}

type UserLoginRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
	Result   LoginResult        `json:"result"`
	UserInfo s2usercrud.UserDao `json:"UserDao"`
}

type LoginResult struct {
	Success         bool  `json:"success"`
	Error           error `json:"error"`
	UserNameMatches bool  `json:"username_matches"`
	PasswordValid   bool  `json:"password_valid"`
	UserEnabled     bool  `json:"enabled"`
}

type SmbPlusPlusClaim struct {
	*jwt.RegisteredClaims
	UserInfo interface{}
}

type TokenRefreshReq struct {
	RefreshToken string `json:"refreshToken"`
}
