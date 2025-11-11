ALTER TABLE deliverys
  ADD CONSTRAINT fk_delivery_order
  FOREIGN KEY (order_id) REFERENCES orders (id);