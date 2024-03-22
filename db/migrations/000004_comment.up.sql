CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS comments (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    post_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    value TEXT NOT NULL,
    created_at timestamptz not null default current_timestamp,
    updated_at timestamptz not null default current_timestamp
);