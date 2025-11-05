-- name: CreateWalletOperation :one
INSERT INTO wallets (operationType, amount)
VALUES ($1, $2)
RETURNING *;

-- name: GetWalletAmount :one
SELECT amount
FROM wallets
WHERE valletId = $1;