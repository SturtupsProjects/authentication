CREATE TABLE company
(
    company_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name       VARCHAR(100) NOT NULL,
    website    VARCHAR(100),
    logo       VARCHAR(255),
    created_at TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP        DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE branches
(
    branch_id  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name       VARCHAR(100) NOT NULL,
    address    VARCHAR(255),
    phone      VARCHAR(15),
    company_id UUID         NOT NULL REFERENCES company (company_id) ON DELETE CASCADE,
    created_at TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    deleted_at bigint           default 0,
    UNIQUE(name, company_id)
);

-- Таблица пользователей
create type role as enum ('admin', 'owner', 'worker');

CREATE TABLE users
(
    user_id      UUID      DEFAULT gen_random_uuid() PRIMARY KEY,
    first_name   VARCHAR(50)        NOT NULL,
    last_name    VARCHAR(50)        NOT NULL,
    email        VARCHAR(100) UNIQUE,
    phone_number VARCHAR(15) UNIQUE NOT NULL,
    password     VARCHAR            NOT NULL,
    role         role               NOT NULL,
    company_id   UUID REFERENCES company (company_id),
    branch_id    UUID REFERENCES branches (branch_id), -- New column for branch association
    created_at   TIMESTAMP DEFAULT NOW()
);

create type type as enum ('client', 'suplier');
create type type_client as enum ('street', 'client');

CREATE TABLE clients
(
    id          UUID      DEFAULT gen_random_uuid() PRIMARY KEY,
    full_name   VARCHAR(60) NOT NULL,
    type        type        NOT NULL,
    client_type type_client NOT NULL,
    address     VARCHAR(50),
    phone       VARCHAR(13),
    company_id  UUID        NOT NULL,
    created_at  TIMESTAMP DEFAULT NOW()
);

