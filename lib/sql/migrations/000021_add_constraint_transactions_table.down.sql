ALTER TABLE transactions,
DROP CONSTRAINT fk_transactions_status,
DROP CONSTRAINT fk_transactions_users,
DROP CONSTRAINT fk_transactions_delivery,
DROP CONSTRAINT fk_transactions_payment;