CREATE TABLE size_product (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200),
    product_id INT,
    additional_costs DECIMAL(10,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);