package s2usercrud

import (
	"context"
	"fmt"
	"log/slog"

	smbplusplus_db "github.com/babbage88/smbplusplus/database/smbplusplus_pg"
	"github.com/babbage88/smbplusplus/internal/hashing"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserCrudPgxImpl struct {
	DbConn *pgxpool.Pool `json:"dbConn"`
}

type TemporaryHolder interface {
	GetAllActiveUsersDao() ([]UserDao, error)
	GetAllActiveRoles([]UserRoleDao, error)
	GetAllAppPermissions([]AppPermissionDao, error)
	GetUserByName(username string) (UserDao, error)
	GetUserById(id uuid.UUID) (UserDao, error)
	UpdateUserEmailById(id uuid.UUID, email string)
	VerifyAlterUser(executionUserId uuid.UUID) (bool, error)
	UpdateUserPasswordWithAuth(execUserId uuid.UUID, targetUserId uuid.UUID, newPassword string) error
	EnableUserById(targetUserId uuid.UUID) (UserDao, error)
	DisableUserById(targetUserId uuid.UUID) (UserDao, error)
	SoftDeleteUserById(targetUserId uuid.UUID) (UserDao, error)
	UpdateUserRoleMapping(targetUserId uuid.UUID, roleId uuid.UUID) error
	DisableUserRoleMapping(targetUserId uuid.UUID, roleId uuid.UUID) error
	CreateOrUpdateUserRole(roleName string, roleDescr string) (*UserRoleDao, error)
	CreateOrUpdateAppPermission(name string, desc string) (*AppPermissionDao, error)
	CreateOrUpdateRolePermisssionMapping(roleId uuid.UUID, permId uuid.UUID) (*RolePermissionMappingDao, error)
	EnableRoleById(id uuid.UUID) error
	DisableRoleById(id uuid.UUID) error
	SoftDeleteRoleById(id uuid.UUID) error
}

type UserCRUD interface {
	NewUser(username string, hashed_pw string, email string) (UserDao, error)
	UpdateUserPasswordById(id uuid.UUID, password string) error
}

func (us *UserCrudPgxImpl) UpdateUserPasswordById(id uuid.UUID, password string) error {
	hashed_pw, _ := hashing.HashPassword(password)

	params := &smbplusplus_db.UpdateUserPasswordByIdParams{ID: id, Password: pgtype.Text{String: hashed_pw, Valid: true}}
	queries := smbplusplus_db.New(us.DbConn)
	err := queries.UpdateUserPasswordById(context.Background(), *params)
	if err != nil {
		slog.Error("Error updating user password in database", slog.String("ID", fmt.Sprint(id)), slog.String("Error", err.Error()))
	}
	return err
}

func (us *UserCrudPgxImpl) NewUser(username string, password string, email string) (UserDao, error) {
	hashed_pw, _ := hashing.HashPassword(password)
	var newuser UserDao
	// Set up parameters for the new user
	params := smbplusplus_db.CreateUserParams{
		Username: pgtype.Text{String: username, Valid: true},
		Password: pgtype.Text{String: hashed_pw, Valid: true},
		Email:    pgtype.Text{String: email, Valid: true},
	}

	queries := smbplusplus_db.New(us.DbConn)
	qry, err := queries.CreateUser(context.Background(), params)
	newuser.ParseUserFromDb(qry)
	return newuser, err
}
