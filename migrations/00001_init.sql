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

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
