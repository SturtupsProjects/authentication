-- Таблица пользователей
CREATE TABLE users
(
    user_id      UUID      DEFAULT gen_random_uuid() PRIMARY KEY,
    first_name   VARCHAR(50)         NOT NULL,
    last_name    VARCHAR(50)         NOT NULL,
    email        VARCHAR(100) UNIQUE,
    phone_number VARCHAR(15) UNIQUE NOT NULL,
    password     VARCHAR             NOT NULL,
    role         VARCHAR(20)         NOT NULL,
    created_at   TIMESTAMP DEFAULT NOW()
);

CREATE TABLE clients
(
    id         UUID      DEFAULT gen_random_uuid() PRIMARY KEY,
    full_name  VARCHAR(60) NOT NULL,
    address    VARCHAR(50),
    phone      VARCHAR(13),
    created_at TIMESTAMP DEFAULT NOW()
);
