CREATE TABLE IF NOT EXISTS master_profiles
(
    id           uuid PRIMARY KEY                DEFAULT uuid_generate_v4(),
    display_name TEXT                   NOT NULL,
    job          TEXT                   NOT NULL,
    description  TEXT,
    point        geography(Point, 4326) NOT NULL,
    address      TEXT,
    avatar_id    TEXT,
    user_id      uuid                   NOT NULL UNIQUE REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE,
    is_enabled   BOOLEAN                NOT NULL DEFAULT TRUE,
    created_at   timestamptz            NOT NULL DEFAULT NOW(),
    updated_at   timestamptz            NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS is_enabled_master_profiles_index on master_profiles (is_enabled);
CREATE INDEX IF NOT EXISTS point_master_profiles_index on master_profiles USING GIST (point);
CREATE INDEX IF NOT EXISTS display_name_master_profiles_index on master_profiles USING GIN (display_name gin_trgm_ops);
CREATE INDEX IF NOT EXISTS job_master_profiles_index on master_profiles USING GIN (job gin_trgm_ops);
CREATE INDEX IF NOT EXISTS address_master_profiles_index on master_profiles USING GIN (address gin_trgm_ops);