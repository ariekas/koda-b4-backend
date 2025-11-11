ALTER TABLE image_products
  ADD CONSTRAINT fk_image_product
  FOREIGN KEY (product_id) REFERENCES products (id);

