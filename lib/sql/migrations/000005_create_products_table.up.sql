CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    discounts_id INT,
    name VARCHAR NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    price_discounts DECIMAL (10,2) DEFAULT 0,
    description TEXT,
    stock INT DEFAULT 0,
    is_flashsale BOOLEAN DEFAULT FALSE,
    is_favorite_product BOOLEAN DEFAULT FALSE,
    category_products_id INT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);