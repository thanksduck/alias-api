-- Rules table
CREATE TABLE "rules" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "username" varchar NOT NULL,
  "alias_email" varchar NOT NULL UNIQUE,
  "destination_email" varchar NOT NULL,
  "is_active" boolean NOT NULL DEFAULT true,
  "comment" varchar NOT NULL DEFAULT '',
  "name" varchar NOT NULL DEFAULT '',
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now()
);

-- Destinations table
CREATE TABLE "destinations" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "username" varchar NOT NULL,
  "destination_email" varchar NOT NULL,
  "domain" varchar NOT NULL,
  "cloudflare_destination_id" varchar NOT NULL UNIQUE,
  "is_verified" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now()
);

-- Foreign Keys
ALTER TABLE "rules" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "rules" ADD FOREIGN KEY ("username") REFERENCES "users" ("username") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "destinations" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "destinations" ADD FOREIGN KEY ("username") REFERENCES "users" ("username") ON DELETE CASCADE ON UPDATE CASCADE;

-- Indexes
CREATE INDEX ON "rules" ("user_id");
CREATE INDEX ON "rules" ("username");

CREATE INDEX ON "destinations" ("user_id");
CREATE INDEX ON "destinations" ("username");
CREATE INDEX ON "destinations" ("domain");
