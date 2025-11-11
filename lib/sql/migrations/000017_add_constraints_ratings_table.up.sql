ALTER TABLE ratings
  ADD CONSTRAINT fk_rating_user
  FOREIGN KEY (user_id) REFERENCES users (id),
  ADD CONSTRAINT fk_rating_product
  FOREIGN KEY (product_id) REFERENCES products (id);