-- Migration: 002_create_user_consent_status
-- Description: Estado actual de consentimientos por usuario (1-a-1 con users)

CREATE TABLE user_consent_status (
    user_id                  UUID        PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    consented_privacy_policy BOOLEAN     NOT NULL DEFAULT FALSE,
    consented_data_campaign  BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at               TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at               TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE  user_consent_status                         IS 'Estado actual de consentimientos del usuario';
COMMENT ON COLUMN user_consent_status.consented_privacy_policy IS 'TRUE si el usuario aceptó la política de privacidad';
COMMENT ON COLUMN user_consent_status.consented_data_campaign  IS 'TRUE si el usuario aceptó el uso de datos para campañas';


-- Migration: 002_create_user_consents
-- Description: Historial de consentimientos — auditoría completa de cada acción

CREATE TABLE user_consents (
    id          UUID        PRIMARY KEY DEFAULT uuidv7(),
    user_id     UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type        VARCHAR(50) NOT NULL,
    accepted    BOOLEAN     NOT NULL,
    version     VARCHAR(20) NOT NULL,
    ip_address  VARCHAR(45) NOT NULl,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_user_consents_user_id ON user_consents (user_id);

COMMENT ON TABLE  user_consents            IS 'Historial de consentimientos del usuario para auditoría';
COMMENT ON COLUMN user_consents.type       IS 'Tipo de consentimiento: privacy_policy | data_campaign';
COMMENT ON COLUMN user_consents.accepted   IS 'TRUE si aceptó, FALSE si rechazó o retiró el consentimiento';
COMMENT ON COLUMN user_consents.version    IS 'Versión del documento legal aceptado (ej: v1.0)';
COMMENT ON COLUMN user_consents.ip_address IS 'IP desde donde se registró el consentimiento';