CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS currencies
(
    id    UUID       NOT NULL DEFAULT gen_random_uuid(),
    date  DATE       NOT NULL,
    base  VARCHAR(3) NOT NULL,
    quote VARCHAR(3) NOT NULL,
    rate  NUMERIC    NOT NULL,

    CONSTRAINT pk_currencies PRIMARY KEY (id),
    CONSTRAINT uq_currency_date UNIQUE (date, base, quote)
);

CREATE INDEX idx_currencies_quote ON currencies (quote);

CREATE INDEX idx_currencies_date ON currencies (date);