CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    users_id INT,
    deliverys_id INT,
    payment_methods_id INT,
    status_transactions_id INT,
    name_user VARCHAR NOT NULL,
    address_user TEXT NOT NULL,
    phone_user VARCHAR,
    email_user VARCHAR,
    total DECIMAL(10,2) NOT NULL,
    invoice_num VARCHAR,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);