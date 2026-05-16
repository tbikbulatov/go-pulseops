-- +goose Up
CREATE TABLE IF NOT EXISTS incidents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    service TEXT NOT NULL,
    environment TEXT NOT NULL,
    severity TEXT NOT NULL,
    status TEXT NOT NULL,
    dedup_key TEXT NOT NULL,
    alert_count INTEGER NOT NULL DEFAULT 1,
    first_seen_at TIMESTAMPTZ NOT NULL,
    last_seen_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT incidents_severity_check CHECK (severity IN ('info', 'warning', 'critical')),
    CONSTRAINT incidents_status_check CHECK (status IN ('open', 'acknowledged', 'resolved')),
    CONSTRAINT incidents_alert_count_check CHECK (alert_count > 0)
);

CREATE UNIQUE INDEX IF NOT EXISTS incidents_active_dedup_idx
    ON incidents (service, environment, dedup_key)
    WHERE status IN ('open', 'acknowledged');

CREATE INDEX IF NOT EXISTS incidents_status_updated_at_idx
    ON incidents (status, updated_at DESC);

CREATE TABLE IF NOT EXISTS incident_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    incident_id UUID NOT NULL REFERENCES incidents(id),
    type TEXT NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS incident_events_incident_id_created_at_idx
    ON incident_events (incident_id, created_at);

CREATE TABLE IF NOT EXISTS processed_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    consumer_name TEXT NOT NULL,
    message_id UUID NOT NULL,
    processed_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT processed_messages_consumer_message_unique UNIQUE (consumer_name, message_id)
);

-- +goose Down
DROP TABLE IF EXISTS processed_messages;

DROP INDEX IF EXISTS incident_events_incident_id_created_at_idx;
DROP TABLE IF EXISTS incident_events;

DROP INDEX IF EXISTS incidents_status_updated_at_idx;
DROP INDEX IF EXISTS incidents_active_dedup_idx;
DROP TABLE IF EXISTS incidents;
