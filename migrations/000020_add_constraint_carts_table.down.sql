ALTER TABLE carts,
DROP CONSTRAINT fk_carts_users,
DROP CONSTRAINT fk_carts_products,
DROP CONSTRAINT fk_carts_size,
DROP CONSTRAINT fk_carts_variant;