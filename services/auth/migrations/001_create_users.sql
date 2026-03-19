-- Migration: 001_create_users
-- Description: Tabla principal de usuarios del sistema de autenticación

CREATE EXTENSION IF NOT EXISTS "pgcrypto"; -- para gen_random_uuid()

CREATE TABLE users (
    id            UUID        PRIMARY KEY DEFAULT uuidv7(),
    email         VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255),                    -- NULL si el proveedor es Google
    first_name     VARCHAR(255) NOT NULL,
    last_name      VARCHAR(255) NOT NULL,
    document_type  VARCHAR(25) NOT NULL,
    document_number VARCHAR(25) NOT NULL,
    phone_number   VARCHAR(50),
    provider      VARCHAR(50)  NOT NULL DEFAULT 'local', -- 'local' | 'google'
    provider_id   VARCHAR(255),                    -- Google sub ID, NULL si es local
    is_verified   BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- Email único globalmente — no puede haber dos cuentas con el mismo email
-- independientemente del proveedor
CREATE UNIQUE INDEX idx_users_email ON users (email);

-- Búsqueda rápida por proveedor (Google OAuth)
CREATE UNIQUE INDEX idx_users_provider ON users (provider, provider_id)
    WHERE provider_id IS NOT NULL;

CREATE UNIQUE INDEX idx_users_document ON users (document_number)

-- Comentarios de columnas para documentación en la DB
COMMENT ON TABLE  users                IS 'Usuarios registrados en el sistema';
COMMENT ON COLUMN users.document_type   IS 'Tipo de documento: DNI | RUC | CE';
COMMENT ON COLUMN users.document_number IS 'Documento de identidad del usuario';
COMMENT ON COLUMN users.provider       IS 'Origen del registro: local | google';
COMMENT ON COLUMN users.provider_id    IS 'ID único del usuario en el proveedor externo (ej. Google sub)';
COMMENT ON COLUMN users.password_hash  IS 'Hash bcrypt de la contraseña. NULL para usuarios OAuth';
COMMENT ON COLUMN users.is_verified    IS 'TRUE si el email fue confirmado o si el registro fue vía OAuth';
