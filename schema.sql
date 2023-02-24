CREATE TABLE IF NOT EXISTS users
(
    id           TEXT                     NOT NULL PRIMARY KEY,
    kind         TEXT                     NOT NULL,
    email        TEXT                     NOT NULL UNIQUE,
    username     TEXT                     NOT NULL UNIQUE,
    pwd_hash     TEXT,
    user_data    jsonb                    NOT NULL,
    created_at   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp,
    updated_at   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp,
    verified_at  TIMESTAMP WITH TIME ZONE          DEFAULT NULL,
    verify_token TEXT                              DEFAULT NULL,
    attributes   jsonb
);
CREATE INDEX IF NOT EXISTS idx_users_kind ON users (kind);

CREATE TABLE IF NOT EXISTS user_keys
(
    key     TEXT NOT NULL PRIMARY KEY,
    user_id TEXT NOT NULL,
    attribs jsonb,

    FOREIGN KEY (user_id) REFERENCES users (id)
);
CREATE INDEX IF NOT EXISTS idx_user_keys_user_id ON user_keys (user_id);