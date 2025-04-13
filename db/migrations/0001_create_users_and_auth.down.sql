-- Remove foreign keys
ALTER TABLE "user_auth" DROP CONSTRAINT IF EXISTS "user_auth_user_id_fkey";
ALTER TABLE "user_auth" DROP CONSTRAINT IF EXISTS "user_auth_username_fkey";
ALTER TABLE "social_profiles" DROP CONSTRAINT IF EXISTS "social_profiles_user_id_fkey";
ALTER TABLE "social_profiles" DROP CONSTRAINT IF EXISTS "social_profiles_username_fkey";

-- Remove indexes
DROP INDEX IF EXISTS "users_username_idx";
DROP INDEX IF EXISTS "users_email_idx";
DROP INDEX IF EXISTS "user_auth_user_id_idx";
DROP INDEX IF EXISTS "social_profiles_user_id_idx";
DROP INDEX IF EXISTS "social_profiles_username_idx";

-- Drop tables
DROP TABLE IF EXISTS "social_profiles";
DROP TABLE IF EXISTS "user_auth";
DROP TABLE IF EXISTS "users";