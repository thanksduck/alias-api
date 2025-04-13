-- First, drop the foreign key constraints
ALTER TABLE "subscriptions" DROP CONSTRAINT IF EXISTS "subscriptions_credit_id_fkey";
ALTER TABLE "subscriptions" DROP CONSTRAINT IF EXISTS "subscriptions_user_id_fkey";
ALTER TABLE "credits" DROP CONSTRAINT IF EXISTS "credits_user_id_fkey";
ALTER TABLE "payments" DROP CONSTRAINT IF EXISTS "payments_user_id_fkey";

-- Drop indexes
DROP INDEX IF EXISTS "subscriptions_status_idx";
DROP INDEX IF EXISTS "subscriptions_user_id_idx";
DROP INDEX IF EXISTS "credits_user_id_idx";
DROP INDEX IF EXISTS "payments_status_idx";
DROP INDEX IF EXISTS "payments_user_id_idx";

-- Finally, drop the tables
DROP TABLE IF EXISTS "subscriptions";
DROP TABLE IF EXISTS "credits";
DROP TABLE IF EXISTS "payments";