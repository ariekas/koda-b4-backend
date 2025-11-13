CREATE TABLE carts (
    id SERIAL PRIMARY KEY,
    users_id INT,
    products_id INT,
    size_products_id INT,
    variant_products_id INT,
    quantity INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);