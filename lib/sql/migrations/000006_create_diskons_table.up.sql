
CREATE TABLE diskons (
  id SERIAL PRIMARY KEY,
  percentage DECIMAL(10,2),
  name VARCHAR(100),
  start_date DATE,
  end_date DATE,
  isActive BOOLEAN DEFAULT false,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);