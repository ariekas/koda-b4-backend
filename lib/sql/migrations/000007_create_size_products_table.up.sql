CREATE TABLE size_products (
    id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL,
    additional_costs DECIMAL(10,2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);