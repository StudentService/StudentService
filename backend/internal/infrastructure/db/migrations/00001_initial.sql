-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Только users без внешних ключей на группы
CREATE TABLE users (
                       id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                       username VARCHAR(255) UNIQUE NOT NULL,
                       email VARCHAR(255) UNIQUE NOT NULL,
                       password_hash VARCHAR(255) NOT NULL,
                       role VARCHAR(50) NOT NULL,
                       first_name VARCHAR(255),
                       last_name VARCHAR(255),
                       created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Индексы
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);

-- Добавим тестового пользователя (пароль: "password123")
INSERT INTO users (username, email, password_hash, role, first_name, last_name)
VALUES (
           'john.doe',
           'john@example.com',
           '$2a$10$YourHashedPasswordHere', -- замените на реальный хеш
           'student',
           'John',
           'Doe'
       );
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS users CASCADE;
DROP EXTENSION IF EXISTS "uuid-ossp";