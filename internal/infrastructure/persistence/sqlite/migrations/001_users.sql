CREATE TABLE IF NOT EXISTS users (
    id TEXT,
    UNIQUE (id)
);

CREATE INDEX IF NOT EXISTS idx_users_id ON users(id);
