CREATE TABLE IF NOT EXISTS accounts (
    id SERIAL PRIMARY KEY,
    user_id UUID,
    username TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
)