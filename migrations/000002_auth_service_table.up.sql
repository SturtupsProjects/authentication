-- Создаем ENUM-типы для статусов и типов транзакций
CREATE TYPE account_status AS ENUM ('active', 'blocked');
CREATE TYPE transaction_category AS ENUM ('deposit', 'charge');

-- Основная таблица баланса компании
CREATE TABLE company_account_balance
(
    id                UUID                    DEFAULT gen_random_uuid() PRIMARY KEY,
    company_id        UUID           NOT NULL REFERENCES company (company_id),
    monthly_fee       DECIMAL(15, 2) NOT NULL,                  -- Месячная оплата
    balance           DECIMAL(15, 2) NOT NULL DEFAULT 0,        -- Баланс счета
    status            account_status NOT NULL DEFAULT 'active', -- Статус аккаунта
    last_payment_date TIMESTAMP,                                -- Дата последнего успешного платежа
    next_due_date     TIMESTAMP      NOT NULL                   -- Дата следующего списания
);

-- Таблица транзакций баланса компании
CREATE TABLE company_balance_transaction
(
    id               UUID                          DEFAULT gen_random_uuid() PRIMARY KEY,
    company_id       UUID                 NOT NULL REFERENCES company (company_id),
    transaction_date TIMESTAMP            NOT NULL DEFAULT NOW(),
    category         transaction_category NOT NULL, -- Тип транзакции (пополнение/списание)
    amount           DECIMAL(15, 2)       NOT NULL CHECK (amount > 0)
);
