CREATE TABLE IF NOT EXISTS identities (
    user_id TEXT PRIMARY KEY,
    provider TEXT NOT NULL,
    external_id TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_identities_userid ON identities(user_id);
CREATE INDEX IF NOT EXISTS idx_identities_provider ON identities(provider);
CREATE INDEX IF NOT EXISTS idx_identities_externalid ON identities(external_id);
