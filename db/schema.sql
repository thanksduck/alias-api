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
  "password" varchar,
  "provider" varchar NOT NULL CHECK (provider IN ('local', 'github', 'google', 'facebook')) DEFAULT 'local',
  "avatar" varchar NOT NULL DEFAULT 'https://n3y.in/ETzABq',
  "password_changed_at" timestamptz NOT NULL DEFAULT '2000-01-01 00:00:00',
  "active" boolean NOT NULL DEFAULT true,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

-- Rules table
CREATE TABLE "rules" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "username" varchar NOT NULL,
  "alias_email" varchar NOT NULL UNIQUE,
  "destination_email" varchar NOT NULL,
  "active" boolean NOT NULL DEFAULT true,
  "comment" varchar,
  "name" varchar,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

-- User auth table
CREATE TABLE "user_auth" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL UNIQUE,
  "username" varchar NOT NULL,
  "password_reset_token" varchar,
  "password_reset_expires" timestamptz,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

-- Destinations table
CREATE TABLE "destinations" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "username" varchar NOT NULL,
  "destination_email" varchar NOT NULL,
  "domain" varchar NOT NULL,
  "cloudflare_destination_id" varchar NOT NULL UNIQUE,
  "verified" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

-- Social profiles table
CREATE TABLE "social_profiles" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "username" varchar NOT NULL,
  "github" varchar,
  "google" varchar,
  "facebook" varchar,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

-- Payments table
CREATE TABLE "payments" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "type" varchar NOT NULL DEFAULT 'credit',
  "gateway" varchar NOT NULL DEFAULT 'phonepe',
  "txn_id" text NOT NULL UNIQUE,
  "amount" bigint NOT NULL,
  "status" varchar NOT NULL CHECK (status IN ('success', 'pending', 'failed')) DEFAULT 'pending',
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

-- Credits table
CREATE TABLE "credits" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "balance" bigint NOT NULL,
  "is_active" boolean NOT NULL DEFAULT true,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

-- Subscriptions table
CREATE TABLE "subscriptions" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "credit_id" bigint,
  "plan" varchar NOT NULL CHECK (plan IN ('star', 'free', 'galaxy')) DEFAULT 'free',
  "price" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "expires_at" timestamptz NOT NULL,
  "status" varchar NOT NULL CHECK (status IN ('active', 'paused', 'cancelled')) DEFAULT 'active'
);

-- Add foreign keys
ALTER TABLE "rules" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "rules" ADD FOREIGN KEY ("username") REFERENCES "users" ("username") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "user_auth" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "user_auth" ADD FOREIGN KEY ("username") REFERENCES "users" ("username") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "destinations" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "destinations" ADD FOREIGN KEY ("username") REFERENCES "users" ("username") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "social_profiles" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "social_profiles" ADD FOREIGN KEY ("username") REFERENCES "users" ("username") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "payments" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "credits" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "subscriptions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "subscriptions" ADD FOREIGN KEY ("credit_id") REFERENCES "credits" ("id") ON DELETE SET NULL ON UPDATE CASCADE;

-- Create indices
CREATE INDEX ON "users" ("username");
CREATE INDEX ON "users" ("email");

CREATE INDEX ON "rules" ("user_id");
CREATE INDEX ON "rules" ("username");

CREATE INDEX ON "user_auth" ("user_id");
CREATE INDEX ON "user_auth" ("username");

CREATE INDEX ON "destinations" ("user_id");
CREATE INDEX ON "destinations" ("username");
CREATE INDEX ON "destinations" ("domain");

CREATE INDEX ON "social_profiles" ("user_id");
CREATE INDEX ON "social_profiles" ("username");

CREATE INDEX ON "payments" ("user_id");
CREATE INDEX ON "payments" ("status");

CREATE INDEX ON "credits" ("user_id");
CREATE INDEX ON "credits" ("balance");

CREATE INDEX ON "subscriptions" ("user_id");
CREATE INDEX ON "subscriptions" ("status");

-- Add comments
COMMENT ON TABLE "users" IS 'Store user accounts information';
COMMENT ON TABLE "rules" IS 'Email forwarding rules for users';
COMMENT ON TABLE "destinations" IS 'Verified email destinations for forwarding';
COMMENT ON TABLE "social_profiles" IS 'Connected social accounts for authentication';
COMMENT ON TABLE "payments" IS 'Payment transactions';
COMMENT ON TABLE "credits" IS 'User credit balance for subscription payments';
COMMENT ON TABLE "subscriptions" IS 'User subscription plans';

COMMENT ON COLUMN "users"."password" IS 'Hashed password for local authentication';
COMMENT ON COLUMN "users"."provider" IS 'Authentication provider (local, github, google, facebook)';

COMMENT ON COLUMN "payments"."amount" IS 'Amount in smallest currency unit (e.g., cents)';
COMMENT ON COLUMN "credits"."balance" IS 'Current credit balance in smallest currency unit';
COMMENT ON COLUMN "subscriptions"."expires_at" IS 'When the current subscription period ends';