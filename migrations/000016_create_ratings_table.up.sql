CREATE TABLE ratings (
    id SERIAL PRIMARY KEY,
    rating INT,
    review TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
