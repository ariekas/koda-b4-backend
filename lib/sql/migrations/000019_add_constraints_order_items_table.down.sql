-- Active: 1762853226689@@127.0.0.1@5432@dbcoffeshop
ALTER TABLE order_items DROP CONSTRAINT fk_order_items_order;
ALTER TABLE order_items DROP CONSTRAINT fk_order_items_product;