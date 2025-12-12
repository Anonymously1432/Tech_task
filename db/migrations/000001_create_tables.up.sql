BEGIN;

CREATE TABLE IF NOT EXISTS wallets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    amount REAL NOT NULL CHECK (amount > 0)
);

CREATE TABLE IF NOT EXISTS requests_history (
    operationType VARCHAR(10) NOT NULL,
    wallet_id UUID NOT NULL,
    amount REAL NOT NULL CHECK (amount > 0),

    CONSTRAINT fk_wallet FOREIGN KEY (wallet_id) REFERENCES wallets(id) ON DELETE CASCADE
);

INSERT INTO wallets (id, amount)
VALUES
    ('f6f9e588-5242-4090-9994-532dd0cfdd4f', '1000'),
    ('ca548f59-5c95-4b59-a3eb-016970585701', '200000'),
    ('862a40a0-9f94-45a9-815f-0fd8473b41f1', '50000');

END;