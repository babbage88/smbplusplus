-- name: CreateUser :one
INSERT INTO users (
    username,
    password,
    email
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetUserById :one
SELECT
    "id",
    "username",
    "password",
    "email",
    "roles",
    "role_ids",
    "created_at",
    "last_modified",
    "enabled",
    "is_deleted"
FROM public.users_with_roles uwr
WHERE "id" = $1;

-- name: GetUserByName :one
SELECT
    "id",
    "username",
    "password",
    "email",
    "roles",
    "role_ids",
    "created_at",
    "last_modified",
    "enabled",
    "is_deleted"
FROM public.users_with_roles uwr
WHERE username = $1;

-- name: GetUserLogin :one
SELECT id, username, "password" , email, "enabled", "roles", "role_ids" FROM public.users_with_roles uwr
WHERE username = $1
LIMIT 1;

-- name: GetUserIdByName :one
SELECT
	id
FROM public.users
where username = $1;

-- name: GetUserNameById :one
SELECT
  "username"
FROM public.users
WHERE id = $1;

-- name: UpdateUserPasswordById :exec
UPDATE users
  set password = $2
WHERE id = $1;

-- name: UpdateUserEmailById :one
UPDATE users
  set email = $2
WHERE id = $1
RETURNING *;

-- name: DeleteUserById :exec
DELETE FROM users
WHERE id = $1;

-- name: SoftDeleteUserById :one
UPDATE users
  set is_deleted = TRUE,
  "enabled" = FALSE
WHERE id = $1
RETURNING *;

-- name: DisableUserById :one
UPDATE users
  set "enabled" = $2
WHERE id = $1
RETURNING *;

-- name: EnableUserById :one
UPDATE users
  set "enabled" = $2
WHERE id = $1
RETURNING *;

-- name: GetAllActiveUsers :many
SELECT
    "id",
    "username",
    "password",
     "email",
    "roles",
    "role_ids",
    "created_at",
    "last_modified",
    "enabled",
    "is_deleted"
FROM public.users_with_roles uwr;

---- name: GetAllUserPermissions :many
SELECT
  "UserId",
  "Username",
  "PermissionId",
  "Permission",
  "Role",
  "LastModified"
FROM
    public.user_permissions_view upv
ORDER BY "UserId" ASC;

-- name: GetUserPermissionsById :many
SELECT
  "UserId",
  "Username",
  "PermissionId",
  "Permission",
  "Role",
  "LastModified"
FROM
    public.user_permissions_view upv
WHERE "UserId" = $1;
--
-- name: VerifyUserPermissionById :one
SELECT EXISTS (
  SELECT
    "UserId",
    "Username",
    "PermissionId",
    "Permission",
    "Role",
    "LastModified"
  FROM
      public.user_permissions_view upv
  WHERE "UserId" = $1 and "Permission" = $2
);

-- name: VerifyUserPermissionByRoleId :one
SELECT EXISTS (
  SELECT
    "RoleId",
    "Role",
    "PermissionId",
    "Permission",
    "Role"
  FROM
      public.role_permissions_view rpv
  WHERE "RoleId" = $1 and "Permission" = $2
);

-- name: InsertOrUpdateUserRoleMappingById :one
INSERT INTO public.user_role_mapping(user_id, role_id, enabled)
VALUES ($1, $2, TRUE)
ON CONFLICT (user_id, role_id)
DO UPDATE SET enabled = TRUE
RETURNING *;

-- name: DisableUserRoleMappingById :one
UPDATE
  public.user_role_mapping
SET
  enabled = FALSE
WHERE user_id = $1 AND role_id = $2
RETURNING *;

-- name: GetRoleIdByName :one
SELECT
  "id" AS "RoleId"
FROM
  public. public.user_roles
WHERE "role_name" = $1;

-- name: InsertOrUpdateUserRole :one
INSERT INTO user_roles (id, role_name, role_description, created_at, last_modified, "enabled", "is_deleted")
VALUES(nextval('user_roles_id_seq'::regclass), $1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, TRUE, false)
ON CONFLICT (role_name)
DO UPDATE SET
	role_description = EXCLUDED.role_description,
	last_modified = CURRENT_TIMESTAMP,
	"enabled" = TRUE,
	"is_deleted" = FALSE
RETURNING *;

-- name: EnableUserRoleById :exec
UPDATE user_roles SET "enabled" = TRUE
WHERE id = $1;

-- name: DisableUserRoleById :exec
UPDATE user_roles SET "enabled" = FALSE
WHERE id = $1;

-- name: SoftDeleteUserRoleById :exec
UPDATE user_roles
SET
"is_deleted" = TRUE,
"enabled" = FALSE
WHERE id = $1;

-- name: HardDeleteUserRoleById :exec
DELETE FROM user_roles
WHERE id = $1;

-- name: InsertOrUpdateAppPermission :one
INSERT INTO app_permissions(id, permission_name, permission_description)
VALUES(nextval('app_permissions_id_seq'::regclass), $1, $2)
ON CONFLICT (permission_name)
DO UPDATE SET
	permission_description = EXCLUDED.permission_description
RETURNING *;

-- name: InsertOrUpdateRolePermissionMapping :one
INSERT INTO role_permission_mapping(id, role_id, permission_id, "enabled", created_at, last_modified)
VALUES(nextval('role_permission_mapping_id_seq'::regclass), $1, $2, true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT(role_id, permission_id)
DO UPDATE SET
  role_id = EXCLUDED.role_id,
  permission_id = EXCLUDED.permission_id,
  "enabled" = true,
  created_at = CURRENT_TIMESTAMP,
  last_modified = CURRENT_TIMESTAMP
RETURNING *;

-- name: GetAllUserRoles :many
SELECT "RoleId", "RoleName", "RoleDescription", "CreatedAt", "LastModified", "Enabled", "IsDeleted"
FROM public.user_roles_active;

-- name: GetAllAppPermissions :many
SELECT id, permission_name, permission_description
FROM public.app_permissions;

-- name: DbHealthCheckRead :one
SELECT id, status, check_type
FROM public.health_check WHERE check_type = 'Read'
LIMIT 1;