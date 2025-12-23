CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_active ON users(active);

-- Insert default admin user (password: admin123)
INSERT INTO users (id, name, email, password_hash, active)
SELECT
    gen_random_uuid(),
    'Admin',
    'admin@bookhub.com',
    '$2a$10$otJUHlZifNL133mJxahlJuDq7w5xv1S3RDeYVCHgikFZ0FOtov2f6',
    TRUE
WHERE NOT EXISTS (SELECT 1 FROM users WHERE email = 'admin@bookhub.com');
