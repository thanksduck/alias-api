DROP INDEX IF EXISTS "webauthn_credentials_user_id_idx";
DROP INDEX IF EXISTS "webauthn_credentials_username_idx";
ALTER TABLE "webauthn_credentials" DROP CONSTRAINT IF EXISTS "webauthn_credentials_user_id_fkey";
ALTER TABLE "webauthn_credentials" DROP CONSTRAINT IF EXISTS "webauthn_credentials_username_fkey";
DROP TABLE IF EXISTS "webauthn_credentials";