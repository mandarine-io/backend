CREATE TABLE IF NOT EXISTS banned_tokens
(
    id
    SERIAL
    PRIMARY
    KEY,
    jti
    TEXT
    NOT
    NULL
    UNIQUE,
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
),
    expired_at BIGINT NOT NULL
    );

CREATE INDEX IF NOT EXISTS expired_at_banned_tokens_index on banned_tokens (expired_at);