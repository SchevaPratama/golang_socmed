CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS payments (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    bank_account_id VARCHAR(255) NOT NULL,
    product_id VARCHAR(255) NOT NULL,
    payment_proof_image_url VARCHAR(255) NOT NULL,
    quantity INT NOT NULL
);