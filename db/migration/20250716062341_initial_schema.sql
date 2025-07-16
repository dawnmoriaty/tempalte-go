-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- Kích hoạt extension để dùng hàm tạo UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Dùng chữ thường, kiểu UUID cho ID, và TIMESTAMPTZ cho thời gian
CREATE TABLE "users" (
  "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  "email" varchar UNIQUE NOT NULL,
  "name" varchar NOT NULL,
  "avatar_url" varchar,
  "password" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "deleted_at" timestamptz
);

CREATE TABLE "products" (
  "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  "code" varchar UNIQUE NOT NULL,
  "name" varchar NOT NULL,
  "image_url" varchar,
  "description" text,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "deleted_at" timestamptz
);

CREATE TABLE "orders" (
  "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  "code" varchar NOT NULL,
  "user_id" uuid NOT NULL REFERENCES "users" ("id"),
  "total_price" integer NOT NULL,
  "status" varchar NOT NULL DEFAULT 'new' CHECK ("status" IN ('new', 'progress', 'done', 'cancelled')),
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "deleted_at" timestamptz
);

CREATE TABLE "order_lines" (
  "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  "order_id" uuid NOT NULL REFERENCES "orders" ("id"),
  "product_id" uuid NOT NULL REFERENCES "products" ("id"),
  "quantity" smallint NOT NULL,
  "price" integer NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "deleted_at" timestamptz
);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
DROP TABLE IF EXISTS "order_lines";
DROP TABLE IF EXISTS "orders";
DROP TABLE IF EXISTS "products";
DROP TABLE IF EXISTS "users";