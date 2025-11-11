  ALTER TABLE size_product
  ADD CONSTRAINT fk_size_product
  FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE CASCADE;