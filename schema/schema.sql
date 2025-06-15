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
-- revoked_tokens table to store revoked tokens
CREATE TABLE IF NOT EXISTS revoked_tokens (
    id SERIAL PRIMARY KEY,
    token VARCHAR(512) NOT NULL,
    user_id INT NOT NULL,
    revoked_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expiry_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
    CONSTRAINT unique_token UNIQUE (token, user_id)
);

CREATE INDEX idx_revoked_tokens_token ON revoked_tokens (token);

-- Index để tìm kiếm theo user_id
CREATE INDEX idx_revoked_tokens_user_id ON revoked_tokens (user_id);

-- Index để tối ưu việc xóa token hết hạn
CREATE INDEX idx_revoked_tokens_expiry ON revoked_tokens (expiry_at);