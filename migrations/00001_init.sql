-- +goose Up

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS s3_buckets (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  provider_name TEXT NOT NULL,
  ui_url TEXT NULL,
  admin_url TEXT NULL
)
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS squared_shares (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  s3_bucket_id uuid REFERENCES s3_buckets (id),
  local_path TEXT NOT NULL,
  smb_path TEXT NULL,
  quota_size TEXT NULL
)
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE public.app_permissions (
	id serial4 NOT NULL,
	permission_name varchar(255) NOT NULL,
	permission_description text NULL,
	CONSTRAINT app_permissions_pkey PRIMARY KEY (id),
	CONSTRAINT unique_permission_name UNIQUE (permission_name)
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE public.auth_tokens (
	id serial4 NOT NULL,
	user_id int4 NULL,
	"token" text NULL,
	expiration timestamp NULL,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	last_modified timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	CONSTRAINT auth_tokens_pkey PRIMARY KEY (id),
	CONSTRAINT check_auth_token_id_nonzero CHECK ((id > 0))
);
CREATE INDEX auth_token_idx_created_at ON public.auth_tokens USING btree (created_at);
CREATE INDEX auth_token_idx_userid ON public.auth_tokens USING btree (user_id);
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE public.auth_tokens ADD CONSTRAINT auth_tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);
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
	id serial4 NOT NULL,
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
	user_id int4 NULL,
	username varchar(255) NULL,
	email varchar(255) NULL,
	deleted_at timestamptz DEFAULT now() NULL,
	deleted_by varchar(255) NULL,
	CONSTRAINT users_audit_pkey PRIMARY KEY (audit_id)
);
-- +goose StatementEnd

-- +goose StatementBegin
-- +goose StatementEnd

-- +goose StatementBegin
-- +goose StatementEnd

-- +goose StatementBegin
-- +goose StatementEnd

-- +goose StatementBegin
-- +goose StatementEnd

-- +goose StatementBegin
-- +goose StatementEnd

-- +goose StatementBegin
-- +goose StatementEnd

-- +goose StatementBegin
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
