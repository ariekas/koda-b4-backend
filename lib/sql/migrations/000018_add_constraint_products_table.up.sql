ALTER TABLE products
ADD CONSTRAINT fk_products_discounts FOREIGN KEY(discounts_id) REFERENCES discounts(id),
ADD CONSTRAINT fk_products_category FOREIGN KEY(category_products_id) REFERENCES category_products(id);
