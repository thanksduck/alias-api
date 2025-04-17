-- Users table
CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar NOT NULL UNIQUE CHECK (username ~ '^[a-z][a-z0-9\-_\.]{3,}$'),
  "name" varchar NOT NULL CHECK (LENGTH(name) >= 4),
  "email" varchar NOT NULL UNIQUE,
  "is_email_verified" boolean NOT NULL DEFAULT false,
  "alias_count" bigint NOT NULL DEFAULT 0,
  "destination_count" bigint NOT NULL DEFAULT 0,
  "is_premium" boolean NOT NULL DEFAULT false,
  "password" varchar NOT NULL DEFAULT '',
  "provider" varchar NOT NULL DEFAULT 'local', 
  "avatar" varchar NOT NULL DEFAULT 'https://n3y.in/ETzABq',
  "password_changed_at" timestamptz NOT NULL DEFAULT('0001-01-01 00:00:00Z'),
  "is_active" boolean NOT NULL DEFAULT true,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now()
);

-- User Auth table
CREATE TABLE "user_auth" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL UNIQUE,
  "username" varchar NOT NULL,
  "password_reset_token" varchar NOT NULL DEFAULT '',
  "password_reset_expires" timestamptz NOT NULL DEFAULT (now() + interval '1 hour'),
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now()
);

-- Social Profiles table
CREATE TABLE "social_profiles" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL UNIQUE,
  "username" varchar NOT NULL UNIQUE,
  "github" varchar NOT NULL DEFAULT 'NULL',
  "google" varchar NOT NULL DEFAULT 'NULL',
  "facebook" varchar NOT NULL DEFAULT 'NULL',
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now()
);

-- Foreign Keys
ALTER TABLE "user_auth" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "user_auth" ADD FOREIGN KEY ("username") REFERENCES "users" ("username") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "social_profiles" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "social_profiles" ADD FOREIGN KEY ("username") REFERENCES "users" ("username") ON DELETE CASCADE ON UPDATE CASCADE;

-- Indexes
CREATE INDEX ON "users" ("username");
CREATE INDEX ON "users" ("email");
CREATE INDEX ON "user_auth" ("user_id");
CREATE INDEX ON "social_profiles" ("user_id");
CREATE INDEX ON "social_profiles" ("username");
