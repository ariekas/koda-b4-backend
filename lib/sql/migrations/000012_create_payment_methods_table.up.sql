CREATE TABLE payment_methods (
    id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL,
    image_payment TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);