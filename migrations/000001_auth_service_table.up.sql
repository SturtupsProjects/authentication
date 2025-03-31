create type role as enum ('admin', 'owner', 'worker');
create type type as enum ('client', 'supplier');
create type type_client as enum ('street', 'client');

create type currency_type as enum ('USD', 'UZS');
create type adjustment_type as enum ('bonus', 'penalty');
create type bonuses_type as enum ('active', 'enactive');

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
    UNIQUE (name, company_id)
);

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
    created_at   TIMESTAMP DEFAULT NOW()
);

CREATE TABLE staff_salary
(
    salary_id     UUID      DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id       UUID           NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    currency_code currency_type  NOT NULL,
    salary_amount DECIMAL(15, 2) NOT NULL CHECK (salary_amount >= 0),
    salary_date   DATE           NOT NULL,
    created_at    TIMESTAMP DEFAULT NOW(),
    updated_at    TIMESTAMP DEFAULT NOW(),
    company_id    UUID           NOT NULL
);

CREATE TABLE salary_adjustments
(
    adjustment_id     UUID                     DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id           UUID            NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    currency_code     currency_type   NOT NULL,
    adjustment_type   adjustment_type NOT NULL, -- например, 'BONUS' или 'PENALTY'
    adjustment_amount DECIMAL(15, 2)  NOT NULL CHECK (adjustment_amount >= 0),
    adjustment_date   DATE            NOT NULL,
    is_active         BOOLEAN         NOT NULL DEFAULT TRUE,
    created_at        TIMESTAMP                DEFAULT NOW(),
    updated_at        TIMESTAMP                DEFAULT NOW(),
    company_id        UUID            NOT NULL
);

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


-- INDEX for tables -----------------------------------------

CREATE INDEX idx_company_id ON company (company_id);

CREATE INDEX idx_branches_company_id ON branches (company_id);
CREATE INDEX idx_branches_not_deleted ON branches (deleted_at) WHERE deleted_at = 0;

CREATE UNIQUE INDEX idx_users_phone ON users (phone_number);
CREATE INDEX idx_users_company_id ON users (company_id);
CREATE INDEX idx_user_user_id ON users (user_id);

CREATE INDEX idx_clients_company_id ON clients (company_id);
CREATE INDEX idx_clients_type ON clients (type);
CREATE INDEX idx_clients_phone ON clients (phone);

CREATE INDEX idx_staff_salary_user_id ON staff_salary (user_id);
CREATE INDEX idx_staff_salary_user_date ON staff_salary (user_id, salary_date DESC);

CREATE INDEX idx_salary_adjustments_user_id ON salary_adjustments (user_id);
CREATE INDEX idx_salary_adjustments_user_date ON salary_adjustments (user_id, adjustment_date DESC);
CREATE INDEX idx_salary_adjustments_user_type_active_date
    ON salary_adjustments (user_id, adjustment_type, is_active, adjustment_date DESC);

