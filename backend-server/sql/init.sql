CREATE TABLE IF NOT EXISTS status (
    id SERIAL PRIMARY KEY,
    designation VARCHAR(255) NOT NULL,
    state VARCHAR(5) CHECK (state IN ('green', 'red')) NOT NULL
);
-- Création de la fonction de déclenchement
CREATE OR REPLACE FUNCTION notify_status_update()
RETURNS TRIGGER AS $$
DECLARE
    payload JSON;
BEGIN
    IF TG_OP = 'INSERT' OR TG_OP = 'UPDATE' THEN
        payload := row_to_json(NEW);
    ELSE
        payload := row_to_json(OLD);
    END IF;

    PERFORM pg_notify('status_update', json_build_object(
        'operation', TG_OP,
        'table', TG_TABLE_NAME,
        'row', payload
    )::text);

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Création du déclencheur pour la table "status"
CREATE TRIGGER status_update_notify
AFTER INSERT OR UPDATE OR DELETE ON status
    FOR EACH ROW EXECUTE FUNCTION notify_status_update();
