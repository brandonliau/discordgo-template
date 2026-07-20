CREATE TABLE IF NOT EXISTS saved_locations (
    user_id TEXT,
    zip     TEXT,
    UNIQUE(user_id, zip)
);

CREATE INDEX IF NOT EXISTS idx_saved_locations_user_id ON saved_locations(user_id);
