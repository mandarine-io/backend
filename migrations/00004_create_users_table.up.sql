CREATE TABLE IF NOT EXISTS users
(
    id                uuid PRIMARY KEY     DEFAULT uuid_generate_v4(),
    username          TEXT        NOT NULL UNIQUE,
    email             TEXT        NOT NULL UNIQUE,
    password          TEXT        NOT NULL,
    role_id           INTEGER     NOT NULL REFERENCES roles (id) ON DELETE RESTRICT ON UPDATE CASCADE,
    is_enabled        BOOLEAN     NOT NULL DEFAULT TRUE,
    is_email_verified BOOLEAN     NOT NULL DEFAULT FALSE,
    is_password_temp  BOOLEAN     NOT NULL DEFAULT TRUE,
    created_at        timestamptz NOT NULL DEFAULT NOW(),
    updated_at        timestamptz NOT NULL DEFAULT NOW(),
    deleted_at        timestamptz
);

CREATE INDEX IF NOT EXISTS is_enabled_users_index on users (is_enabled);
CREATE INDEX IF NOT EXISTS is_email_verified_users_index on users (is_email_verified);
CREATE INDEX IF NOT EXISTS is_password_temp_users_index on users (is_password_temp);
CREATE INDEX IF NOT EXISTS deleted_at_users_index on users (deleted_at);