package s2auth

type AuthService interface {
	LoginUser(username string, password string) (AuthTokenDao, error)
	VerifyUserPermission(id uuid, permission string) error
}
