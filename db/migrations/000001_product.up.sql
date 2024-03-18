CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS products (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    price BIGINT NOT NULL,
    imageUrl VARCHAR(2048),
    -- Allow longer URLs for images
    stock INT NOT NULL,
    condition VARCHAR(50),
    isPurchasable BOOLEAN DEFAULT TRUE,
    tags TEXT,
    userId VARCHAR(255) NOT NULL
);
CREATE TABLE IF NOT EXISTS product_tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    product_id uuid NOT NULL
);