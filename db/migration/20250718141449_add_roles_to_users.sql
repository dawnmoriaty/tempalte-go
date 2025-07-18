-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- Tạo bảng roles
CREATE TABLE "roles" (
  "id" serial PRIMARY KEY,
  "name" varchar UNIQUE NOT NULL -- Ví dụ: 'admin', 'user'
);

-- Thêm một vài role mặc định
INSERT INTO "roles" (name) VALUES ('admin'), ('user');

-- Thêm cột role_id vào bảng users
-- Mặc định gán cho tất cả user hiện tại là role 'user' (id = 2)
ALTER TABLE "users" ADD COLUMN "role_id" integer NOT NULL DEFAULT 2;

-- Thêm khóa ngoại
ALTER TABLE "users" ADD FOREIGN KEY ("role_id") REFERENCES "roles" ("id");

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
