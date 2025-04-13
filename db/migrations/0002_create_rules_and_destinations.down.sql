-- Drop indexes
DROP INDEX IF EXISTS "destinations_domain_idx";
DROP INDEX IF EXISTS "destinations_username_idx";
DROP INDEX IF EXISTS "destinations_user_id_idx";

DROP INDEX IF EXISTS "rules_username_idx";
DROP INDEX IF EXISTS "rules_user_id_idx";

-- Drop foreign keys
ALTER TABLE "destinations" DROP CONSTRAINT IF EXISTS "destinations_username_fkey";
ALTER TABLE "destinations" DROP CONSTRAINT IF EXISTS "destinations_user_id_fkey";

ALTER TABLE "rules" DROP CONSTRAINT IF EXISTS "rules_username_fkey";
ALTER TABLE "rules" DROP CONSTRAINT IF EXISTS "rules_user_id_fkey";

-- Drop tables
DROP TABLE IF EXISTS "destinations";
DROP TABLE IF EXISTS "rules";