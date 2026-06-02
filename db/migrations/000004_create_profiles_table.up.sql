CREATE TABLE profiles (
user_id INTEGER PRIMARY KEY,
full_name VARCHAR(100),
photo VARCHAR(255),
phone VARCHAR(20) UNIQUE,
created_at TIMESTAMP DEFAULT NOW(),
updated_at TIMESTAMP DEFAULT NOW(),

CONSTRAINT fk_profile_user
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE

);