ALTER TABLE order_items
ADD CONSTRAINT fk_order_item_size
FOREIGN KEY (size_product_id)
REFERENCES size_product(id)
ON DELETE SET NULL;

ALTER TABLE order_items
ADD CONSTRAINT fk_order_item_variant
FOREIGN KEY (variant_id)
REFERENCES variant(id)
ON DELETE SET NULL;
