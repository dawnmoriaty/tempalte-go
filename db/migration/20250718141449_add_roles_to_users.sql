-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- Tạo bảng roles
CREATE TABLE "roles" (
  "id" serial PRIMARY KEY,
  "name" varchar UNIQUE NOT NULL
);

-- Thêm một vài role mặc định
INSERT INTO "roles" (name) VALUES ('ROLE_ADMIN'), ('ROLE_USER');

-- Tạo bảng nối user_roles để có quan hệ Nhiều-Nhiều
CREATE TABLE "user_roles" (
  "user_id" uuid NOT NULL,
  "role_id" integer NOT NULL,
  PRIMARY KEY ("user_id", "role_id")
);

-- Thêm các khóa ngoại
ALTER TABLE "user_roles" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;
ALTER TABLE "user_roles" ADD FOREIGN KEY ("role_id") REFERENCES "roles" ("id") ON DELETE CASCADE;

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
DROP TABLE IF EXISTS "user_roles";
DROP TABLE IF EXISTS "roles";