CREATE TABLE IF NOT EXISTS users (
                                     user_id SERIAL PRIMARY KEY,
                                     email VARCHAR(255) UNIQUE NOT NULL,
                                     password VARCHAR(255) NOT NULL,
                                     name VARCHAR(255) NOT NULL,
                                     phone VARCHAR(15),
                                     role VARCHAR(50) NOT NULL DEFAULT 'user',
                                     created_at TIMESTAMP NOT NULL,
                                     updated_at TIMESTAMP NOT NULL,
                                     last_login TIMESTAMP
);
