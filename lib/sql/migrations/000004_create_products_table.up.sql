
CREATE TABLE products (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100),
  price DECIMAL(10,2),
  description TEXT,
  stock INT,
  isFlashSale BOOLEAN DEFAULT false,
  isFavorite_product BOOLEAN DEFAULT false,
  category_product_id INT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);