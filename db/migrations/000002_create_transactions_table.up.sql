CREATE TABLE transactions (
id SERIAL PRIMARY KEY,
type VARCHAR(20) NOT NULL,
reference_code VARCHAR(50) NOT NULL UNIQUE,
status VARCHAR(20) DEFAULT 'pending',
created_at TIMESTAMP DEFAULT NOW(),
updated_at TIMESTAMP
);