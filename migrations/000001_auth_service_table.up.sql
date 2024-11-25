-- Таблица пользователей
CREATE TABLE users
(
    user_id      UUID      DEFAULT gen_random_uuid() PRIMARY KEY,
    first_name   VARCHAR(50)         NOT NULL,
    last_name    VARCHAR(50)         NOT NULL,
    email        VARCHAR(100) UNIQUE NOT NULL,
    phone_number VARCHAR(15),
    password     VARCHAR             NOT NULL,
    role         VARCHAR(20)         NOT NULL,
    created_at   TIMESTAMP DEFAULT NOW()
);


