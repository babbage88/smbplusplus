-- +goose Up

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.s3_buckets (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  provider_name TEXT NOT NULL,
  ui_url TEXT NULL,
  admin_url TEXT NULL
)
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.squared_shares (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  s3_bucket_id uuid REFERENCES s3_buckets (id),
  local_path TEXT NOT NULL,
  smb_path TEXT NULL,
  quota_size TEXT NULL
)
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE public.app_permissions (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	permission_name varchar(255) NOT NULL,
	permission_description text NULL,
	CONSTRAINT unique_permission_name UNIQUE (permission_name)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE public.health_check (
	id int4 NOT NULL,
	status varchar(255) NULL,
	check_type varchar(255) NULL,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	last_modified timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE public.role_permission_mapping (
	id serial4 NOT NULL,
	role_id int4 NOT NULL,
	permission_id int4 NOT NULL,
	enabled bool DEFAULT true NOT NULL,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	last_modified timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	CONSTRAINT role_permission_mapping_pkey PRIMARY KEY (id),
	CONSTRAINT unique_perm_role_id UNIQUE (permission_id, role_id)
);


-- public.role_permission_mapping foreign keys

ALTER TABLE public.role_permission_mapping ADD CONSTRAINT fk_permission FOREIGN KEY (permission_id) REFERENCES public.app_permissions(id) ON DELETE CASCADE;
ALTER TABLE public.role_permission_mapping ADD CONSTRAINT fk_role FOREIGN KEY (role_id) REFERENCES public.user_roles(id) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE public.user_role_mapping (
	id serial4 NOT NULL,
	user_id int4 NOT NULL,
	role_id int4 NOT NULL,
	enabled bool DEFAULT true NOT NULL,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	last_modified timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	CONSTRAINT unique_user_role_id UNIQUE (user_id, role_id),
	CONSTRAINT user_role_mapping_pkey PRIMARY KEY (id)
);

-- public.user_role_mapping foreign keys
ALTER TABLE public.user_role_mapping ADD CONSTRAINT fk_role FOREIGN KEY (role_id) REFERENCES public.user_roles(id) ON DELETE CASCADE;
ALTER TABLE public.user_role_mapping ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE public.user_roles (
	id serial4 NOT NULL,
	role_name varchar(255) NOT NULL,
	role_description text NULL,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	last_modified timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	enabled bool DEFAULT true NOT NULL,
	is_deleted bool DEFAULT false NOT NULL,
	CONSTRAINT unique_role_name UNIQUE (role_name),
	CONSTRAINT user_roles_pkey PRIMARY KEY (id)
);
CREATE INDEX user_roles_idx_created_at ON public.user_roles USING btree (created_at);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE public.users (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	username varchar(255) NULL,
	"password" text NULL,
	email varchar(255) NULL,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	last_modified timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	enabled bool DEFAULT true NOT NULL,
	is_deleted bool DEFAULT false NOT NULL,
	CONSTRAINT check_usesr_id_nonzero CHECK ((id > 0)),
	CONSTRAINT unique_email UNIQUE (email),
	CONSTRAINT unique_username UNIQUE (username),
	CONSTRAINT users_pkey PRIMARY KEY (id)
);
CREATE INDEX users_idx_created ON public.users USING btree (created_at);
CREATE INDEX users_idx_user_id ON public.users USING btree (id, username);

-- Table Triggers

create trigger user_delete_trigger after
delete
    on
    public.users for each row execute function log_user_deletion();
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE public.users_audit (
	audit_id serial4 NOT NULL,
	user_id uuid REFERENCES users (id) NULL,
	username varchar(255) NULL,
	email varchar(255) NULL,
	deleted_at timestamptz DEFAULT now() NULL,
	deleted_by varchar(255) NULL,
	CONSTRAINT users_audit_pkey PRIMARY KEY (audit_id)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION log_user_deletion()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO users_audit (user_id, username, email, deleted_at, deleted_by)
    VALUES (OLD.id, OLD.username, OLD.email, now(), current_user);
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER user_delete_trigger
AFTER DELETE ON users
FOR EACH ROW
EXECUTE FUNCTION log_user_deletion();
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE public.auth_tokens (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id uuid REFERENCES users (id),
	"token" text NULL,
	expiration timestamp NULL,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	last_modified timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL
);
CREATE INDEX auth_token_idx_created_at ON public.auth_tokens USING btree (created_at);
CREATE INDEX auth_token_idx_userid ON public.auth_tokens USING btree (user_id);
-- +goose StatementEnd

-- +goose StatementBegin
-- public.role_permissions_view source
CREATE OR REPLACE VIEW public.role_permissions_view
AS SELECT ur.id AS "RoleId",
    ur.role_name AS "Role",
    ap.id AS "PermissionId",
    ap.permission_name AS "Permission"
   FROM user_roles ur
     LEFT JOIN role_permission_mapping rpm ON rpm.role_id = ur.id
     LEFT JOIN app_permissions ap ON rpm.permission_id = ap.id
  WHERE ur.enabled = true AND rpm.enabled = true
  GROUP BY ur.id, ur.role_name, ap.id, ap.permission_name
  ORDER BY ur.id, ap.id;
-- +goose StatementEnd

-- +goose StatementBegin
-- public.user_permissions_view source

CREATE OR REPLACE VIEW public.user_permissions_view
AS SELECT u.id AS "UserId",
    u.username AS "Username",
    ap.id AS "PermissionId",
    ap.permission_name AS "Permission",
    ur.role_name AS "Role",
    urm.last_modified AS "LastModified"
   FROM user_role_mapping urm
     LEFT JOIN user_roles ur ON ur.id = urm.role_id
     LEFT JOIN users u ON u.id = urm.user_id
     LEFT JOIN role_permission_mapping rpm ON rpm.role_id = urm.role_id
     LEFT JOIN app_permissions ap ON ap.id = rpm.permission_id
  WHERE ur.enabled = true
  ORDER BY u.id;
-- +goose StatementEnd

-- +goose StatementBegin
-- public.user_roles_active source

CREATE OR REPLACE VIEW public.user_roles_active
AS SELECT user_roles.id AS "RoleId",
    user_roles.role_name AS "RoleName",
    user_roles.role_description AS "RoleDescription",
    user_roles.created_at AS "CreatedAt",
    user_roles.last_modified AS "LastModified",
    user_roles.enabled AS "Enabled",
    user_roles.is_deleted AS "IsDeleted"
   FROM user_roles
  WHERE user_roles.is_deleted IS FALSE;
-- +goose StatementEnd

-- +goose StatementBegin
-- public.users_with_roles source

CREATE OR REPLACE VIEW public.users_with_roles
AS SELECT u.id,
    u.username,
    u.password,
    u.email,
    COALESCE(array_agg(ur.role_name) FILTER (WHERE ur.role_name IS NOT NULL), ARRAY['None'::text]::character varying[]) AS roles,
    COALESCE(array_agg(urm.role_id) FILTER (WHERE urm.role_id IS NOT NULL), ARRAY[0]) AS role_ids,
    u.created_at,
    u.last_modified,
    u.enabled,
    u.is_deleted
   FROM users u
     LEFT JOIN user_role_mapping urm ON u.id = urm.user_id AND urm.enabled = true
     LEFT JOIN user_roles ur ON urm.role_id = ur.id
  WHERE u.is_deleted = false
  GROUP BY u.id, u.username, u.password, u.email, u.created_at, u.last_modified, u.enabled, u.is_deleted;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW role_permissions_view;
DROP VIEW user_permissions_view;
DROP VIEW user_roles_active;
DROP VIEW users_with_roles;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TRIGGER user_delete_trigger ON public.users;
DROP FUNCTION log_user_deletion;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE public.s3_buckets; 
DROP TABLE public.users_audit;
DROP TABLE public.users;
DROP TABLE public.app_permissions;
DROP TABLE public.auth_tokens;
DROP TABLE public.health_check;
DROP TABLE public.role_permission_mapping;
DROP TABLE public.user_roles;
DROP TABLE public.user_role_mapping;
-- +goose StatementEnd
