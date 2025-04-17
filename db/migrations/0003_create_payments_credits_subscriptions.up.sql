-- Payments table
CREATE TABLE "payments" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "type" varchar NOT NULL DEFAULT 'credit',
  "gateway" varchar NOT NULL DEFAULT 'phonepe',
  "txn_id" text NOT NULL UNIQUE,
  "amount" bigint NOT NULL,
  "status" varchar NOT NULL CHECK (status IN ('success', 'pending', 'failed')) DEFAULT 'pending',
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now()
);

-- Credits table
CREATE TABLE "credits" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL UNIQUE,
  "balance" bigint NOT NULL,
  "is_active" boolean NOT NULL DEFAULT true,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now()
);

-- Subscriptions table
CREATE TABLE "subscriptions" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "credit_id" bigint NOT NULL,
  "plan" varchar NOT NULL CHECK (plan IN ('star', 'free', 'galaxy')) DEFAULT 'free',
  "price" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "expires_at" timestamptz NOT NULL,
  "status" varchar NOT NULL CHECK (status IN ('active', 'paused', 'cancelled')) DEFAULT 'active'
);

-- Foreign Keys
ALTER TABLE "payments" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "credits" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "subscriptions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "subscriptions" ADD FOREIGN KEY ("credit_id") REFERENCES "credits" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- Indexes
CREATE INDEX ON "payments" ("user_id");
CREATE INDEX ON "payments" ("status");

CREATE INDEX ON "credits" ("user_id");

CREATE INDEX ON "subscriptions" ("user_id");
CREATE INDEX ON "subscriptions" ("status");