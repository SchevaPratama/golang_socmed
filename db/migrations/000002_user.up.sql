CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS users (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NULL,
    phone VARCHAR(255) NULL,
    name VARCHAR(255) NOT NULL,
    friends VARCHAR[] NOT NULL default array[]::varchar[],
    password TEXT NOT NULL,
    createdAt timestamptz not null default current_timestamp,
    updatedAt timestamptz not null default current_timestamp
);