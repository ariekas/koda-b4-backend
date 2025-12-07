CREATE TABLE transaction_items (
    id SERIAL PRIMARY KEY,
    transactions_id INT,
    products_id INT,
    quantity INT DEFAULT 1,
    subtotal DECIMAL(10,2) NOT NULL,
    variant_products_id INT,
    size_products_id INT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);