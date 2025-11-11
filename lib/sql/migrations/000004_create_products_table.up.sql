
CREATE TABLE products (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100),
  price NUMERIC(10,2),
  description TEXT,
  product_size VARCHAR(10),
  stock INT,
  isFlashSale BOOLEAN DEFAULT false,
  isFavorite_product BOOLEAN DEFAULT false,
  temperature VARCHAR(50),
  category_product_id INT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);