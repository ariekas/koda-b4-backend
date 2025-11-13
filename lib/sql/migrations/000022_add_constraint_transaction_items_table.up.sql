ALTER TABLE transaction_items
ADD CONSTRAINT fk_transaction_items_transactions FOREIGN KEY (transactions_id) REFERENCES transactions (id),
ADD CONSTRAINT fk_transaction_items_products FOREIGN KEY (products_id) REFERENCES products (id),
ADD CONSTRAINT fk_transaction_items_variant FOREIGN KEY (variant_products_id) REFERENCES variant_products (id),
ADD CONSTRAINT fk_transaction_items_size FOREIGN KEY (size_products_id) REFERENCES size_products (id);