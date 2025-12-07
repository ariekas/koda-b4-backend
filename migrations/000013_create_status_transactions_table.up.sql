CREATE TABLE status_transactions(
    id SERIAL PRIMARY KEY,
    status VARCHAR,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);