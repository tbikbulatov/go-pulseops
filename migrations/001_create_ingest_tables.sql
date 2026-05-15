-- +goose Up
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS integrations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    key TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT integrations_status_check CHECK (status IN ('active', 'inactive'))
);

CREATE TABLE IF NOT EXISTS alerts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    integration_id UUID NOT NULL REFERENCES integrations(id),
    external_id TEXT NOT NULL,
    service TEXT NOT NULL,
    environment TEXT NOT NULL,
    severity TEXT NOT NULL,
    name TEXT NOT NULL,
    message TEXT NOT NULL,
    dedup_key TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT alerts_severity_check CHECK (severity IN ('info', 'warning', 'critical')),
    CONSTRAINT alerts_integration_external_id_unique UNIQUE (integration_id, external_id)
);

CREATE INDEX IF NOT EXISTS alerts_integration_id_idx ON alerts (integration_id);
CREATE INDEX IF NOT EXISTS alerts_dedup_lookup_idx ON alerts (service, environment, dedup_key);

CREATE TABLE IF NOT EXISTS outbox_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    aggregate_type TEXT NOT NULL,
    aggregate_id UUID NOT NULL,
    event_type TEXT NOT NULL,
    payload JSONB NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    attempts INTEGER NOT NULL DEFAULT 0,
    next_attempt_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    published_at TIMESTAMPTZ,
    CONSTRAINT outbox_events_status_check CHECK (status IN ('pending', 'published', 'failed')),
    CONSTRAINT outbox_events_attempts_check CHECK (attempts >= 0)
);

CREATE INDEX IF NOT EXISTS outbox_events_pending_idx
    ON outbox_events (status, next_attempt_at, created_at)
    WHERE status = 'pending';

INSERT INTO integrations (key, name, status)
VALUES ('demo-key', 'Demo Integration', 'active')
ON CONFLICT (key) DO NOTHING;

-- +goose Down
DROP INDEX IF EXISTS outbox_events_pending_idx;
DROP TABLE IF EXISTS outbox_events;

DROP INDEX IF EXISTS alerts_dedup_lookup_idx;
DROP INDEX IF EXISTS alerts_integration_id_idx;
DROP TABLE IF EXISTS alerts;

DROP TABLE IF EXISTS integrations;
