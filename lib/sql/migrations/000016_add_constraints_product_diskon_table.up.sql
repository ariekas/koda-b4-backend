-- Active: 1762853226689@@127.0.0.1@5432@dbcoffeshop
ALTER TABLE product_diskon
  ADD CONSTRAINT fk_product_diskon_product
  FOREIGN KEY (product_id) REFERENCES products (id),
  ADD CONSTRAINT fk_product_diskon_diskon
  FOREIGN KEY (diskon_id) REFERENCES diskons (id);