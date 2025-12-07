ALTER TABLE carts
ADD CONSTRAINT fk_carts_users FOREIGN KEY (users_id) REFERENCES users (id),
ADD CONSTRAINT fk_carts_products FOREIGN KEY (products_id) REFERENCES products (id),
ADD CONSTRAINT fk_carts_size FOREIGN KEY (size_products_id) REFERENCES size_products (id),
ADD CONSTRAINT fk_carts_variant FOREIGN KEY (variant_products_id) REFERENCES variant_products (id);