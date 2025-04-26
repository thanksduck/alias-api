CREATE TABLE "webauthn_credentials" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "username" varchar NOT NULL,
  "credential_id" bytea NOT NULL UNIQUE,
  "public_key" bytea NOT NULL,
  "sign_count" bigint NOT NULL DEFAULT 0,
  "transports" text[] NOT NULL DEFAULT '{}',
  "authenticator_aaguid" varchar NOT NULL DEFAULT '',
  "is_backup" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now()
);

-- Foreign Keys
ALTER TABLE "webauthn_credentials"
ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "webauthn_credentials"
ADD FOREIGN KEY ("username") REFERENCES "users" ("username") ON DELETE CASCADE ON UPDATE CASCADE;

-- Indexes
CREATE INDEX ON "webauthn_credentials" ("user_id");
CREATE INDEX ON "webauthn_credentials" ("username");