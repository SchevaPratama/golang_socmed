CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS posts (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    post_in_html TEXT NOT NULL,
    tags VARCHAR[] NOT NULL default array[]::varchar[],
    user_id TEXT NOT NULL,
    created_at timestamptz not null default current_timestamp,
    updated_at timestamptz not null default current_timestamp
);