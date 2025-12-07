CREATE TABLE profile (
    id SERIAL PRIMARY KEY,
    pic TEXT,
    phone VARCHAR,
    address TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);