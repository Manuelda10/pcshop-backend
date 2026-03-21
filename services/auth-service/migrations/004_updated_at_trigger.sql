-- Migration: 003_updated_at_trigger
-- Description: Trigger para actualizar updated_at automáticamente en cada UPDATE
-- Esto evita tener que setear updated_at manualmente en cada query de actualización

CREATE OR REPLACE FUNCTION trigger_set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Aplicar trigger a la tabla users
CREATE TRIGGER set_updated_at_users
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_updated_at();