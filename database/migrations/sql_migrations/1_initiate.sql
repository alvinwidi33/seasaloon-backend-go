-- +migrate Up
-- +migrate StatementBegin
CREATE TABLE users (
    id            UUID NOT NULL PRIMARY KEY,
    username      VARCHAR(256),
    password      VARCHAR(256),
    role          VARCHAR(256),
    is_active     BOOLEAN
);

CREATE TABLE customer (
    customer_id   UUID NOT NULL PRIMARY KEY,
    user_id       UUID REFERENCES users(id) ON DELETE SET NULL,
    status        VARCHAR(256)
);

CREATE TABLE admin (
    admin_id   UUID NOT NULL PRIMARY KEY,
    user_id    UUID REFERENCES users(id) ON DELETE SET NULL
);

CREATE TABLE saloon (
    id            SERIAL NOT NULL PRIMARY KEY,
    name          VARCHAR(256),
    location      VARCHAR(256),
    open          TIME,
    close         TIME,
    is_delete     BOOLEAN
);

CREATE TABLE reservation (
    id            SERIAL NOT NULL PRIMARY KEY,
    location      VARCHAR(256),
    services      JSON,
    start         TIME,
    done          TIME,
    is_done       BOOLEAN,
    is_cancel     BOOLEAN,
    rating        BIGINT,
    feedback      VARCHAR(256),
    customer_id   UUID REFERENCES customer(customer_id) ON DELETE SET NULL,
    saloon_id     BIGINT REFERENCES saloon(id) ON DELETE SET NULL
);

-- +migrate StatementEnd