  ALTER TABLE variant
  ADD CONSTRAINT fk_variant_product
  FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE CASCADE;