ALTER Table transactions
ADD CONSTRAINT fk_transactions_status FOREIGN KEY (status_transactions_id) REFERENCES status_transactions (id),
ADD CONSTRAINT fk_transactions_users FOREIGN KEY (users_id) REFERENCES users (id),
ADD CONSTRAINT fk_transactions_delivery FOREIGN KEY (deliverys_id) REFERENCES deliverys (id),
ADD CONSTRAINT fk_transactions_payment FOREIGN KEY (payment_methods_id) REFERENCES payment_methods (id);