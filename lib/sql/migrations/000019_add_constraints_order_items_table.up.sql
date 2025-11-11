ALTER TABLE order_items
  ADD CONSTRAINT fk_order_items_order
  FOREIGN KEY (order_id) REFERENCES orders (id),
  ADD CONSTRAINT fk_order_items_product
  FOREIGN KEY (product_id) REFERENCES products (id);