CREATE TABLE IF NOT EXISTS users (
    id TEXT,
    UNIQUE (id)
);

CREATE INDEX IF NOT EXISTS idx_users_userid ON users(id);
