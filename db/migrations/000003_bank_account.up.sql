CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS bank_accounts (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    number VARCHAR(255) NOT NULL,
    bank_name VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL
);