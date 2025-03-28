package s2auth

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type LoginHandler struct {
	Service AuthService `json:"authService"`
}

func LoginHandleFunc(auth_svc AuthService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Username and Pasword in request body as json.
		// in:body
		var loginReq *UserLoginRequest
		json.NewDecoder(r.Body).Decode(&loginReq)
		LoginResult := auth_svc.Login(loginReq)

		if LoginResult.Result.Success {
			token, err := auth_svc.CreateAuthTokenOnLogin(LoginResult.UserInfo.Id, LoginResult.UserInfo.RoleIds, LoginResult.UserInfo.Email)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				slog.Error("Error verifying password", slog.String("Error", err.Error()))
			}
			jsonResponse, _ := json.Marshal(token)
			w.WriteHeader(http.StatusOK)
			w.Write(jsonResponse)
			return
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Invalid credentials", LoginResult.Result.Error)
		}
	}
}

// swagger:route POST /login login idOfloginEndpoint
// Login a user and return token.
// responses:
//   200: AuthToken

func (l *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	LoginHandleFunc(l.Service)
}
