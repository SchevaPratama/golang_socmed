CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS users (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NULL,
    phone VARCHAR(255) NULL,
    name VARCHAR(255) NOT NULL,
    password TEXT NOT NULL
);