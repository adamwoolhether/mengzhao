CREATE TABLE IF NOT EXISTS images (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES accounts(user_id),
    batch_id UUID NOT NULL,
    status INTEGER NOT NULL DEFAULT 1,
    prompt TEXT NOT NULL,
    deleted BOOLEAN NOT NULL DEFAULT FALSE,
    img_loc TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
)