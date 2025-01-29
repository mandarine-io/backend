CREATE TABLE IF NOT EXISTS master_services
(
    id                uuid PRIMARY KEY     DEFAULT uuid_generate_v4(),
    name              TEXT        NOT NULL,
    description       TEXT,
    min_price         money,
    max_price         money,
    min_interval      interval,
    max_interval      interval,
    avatar_id         TEXT,
    master_profile_id uuid        NOT NULL REFERENCES master_profiles (user_id),
    created_at        timestamptz NOT NULL DEFAULT NOW(),
    updated_at        timestamptz NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS name_master_services_index on master_services USING GIN (name gin_trgm_ops);
CREATE INDEX IF NOT EXISTS min_price_master_services_index on master_services (min_price);
CREATE INDEX IF NOT EXISTS max_price_master_services_index on master_services (max_price);
CREATE INDEX IF NOT EXISTS min_interval_master_services_index on master_services (min_interval);
CREATE INDEX IF NOT EXISTS max_interval_master_services_index on master_services (max_interval);