CREATE TABLE IF NOT EXISTS roles
(
    id
    SERIAL
    PRIMARY
    KEY,
    name
    TEXT
    NOT
    NULL
    UNIQUE,
    description
    TEXT,
    created_at
    timestamptz
    NOT
    NULL
    DEFAULT
    NOW
(
),
    updated_at timestamptz NOT NULL DEFAULT NOW
(
)
    );