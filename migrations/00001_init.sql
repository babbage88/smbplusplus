-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS squared_shares(
  id uuid PRIMARY KEY DEFAULT gen_random_uuid()
  local_path TEXT NOT NULL
  s3_bucket_id uuiD NULL
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
