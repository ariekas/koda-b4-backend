CREATE TABLE product_images (
    id SERIAL PRIMARY KEY,
    image TEXT NOT NULL,
    products_id INT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);