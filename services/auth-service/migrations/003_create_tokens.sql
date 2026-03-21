-- Migration: 002_create_tokens
-- Description: Tablas para refresh tokens y verificación de email

-- -------------------------
-- Refresh Tokens
-- -------------------------

CREATE TABLE refresh_tokens (
    id         UUID        PRIMARY KEY DEFAULT uuidv7(),
    user_id    UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(64) NOT NULL,   -- SHA-256 hex del token real (64 chars)
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Búsqueda principal: buscar token por su hash al hacer refresh/logout
CREATE UNIQUE INDEX idx_refresh_tokens_hash ON refresh_tokens (token_hash);

-- Útil para invalidar todos los tokens de un usuario (logout de todos los dispositivos)
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens (user_id);

-- Índice para limpiar tokens expirados con un job periódico:
-- DELETE FROM refresh_tokens WHERE expires_at < NOW();
CREATE INDEX idx_refresh_tokens_expires ON refresh_tokens (expires_at);

COMMENT ON TABLE refresh_tokens IS 'Refresh tokens activos por usuario. Se invalidan en logout o al rotar.';
COMMENT ON COLUMN refresh_tokens.token_hash IS 'SHA-256 del token real. Nunca se almacena el token en plano.';

-- -------------------------
-- Email Verifications
-- -------------------------

CREATE TABLE email_verifications (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    code_hash  VARCHAR(255) NOT NULL,  -- bcrypt hash del código de 6 dígitos
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Un usuario solo puede tener un código activo a la vez
    CONSTRAINT uq_email_verifications_user UNIQUE (user_id)
);

CREATE INDEX idx_email_verifications_expires ON email_verifications (expires_at);

COMMENT ON TABLE  email_verifications           IS 'Códigos de verificación de email pendientes de confirmar.';
COMMENT ON COLUMN email_verifications.code_hash IS 'bcrypt hash del código de 6 dígitos enviado al usuario.';