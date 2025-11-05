BEGIN;

CREATE TABLE IF NOT EXISTS wallets (
    valletId UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    operationType VARCHAR(10) NOT NULL,
    amount INTEGER NOT NULL
);

END;