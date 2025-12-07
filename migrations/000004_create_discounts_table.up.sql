CREATE TABLE discounts(
    id SERIAL PRIMARY KEY,
    name VARCHAR,
    diskon DECIMAL(10,2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);