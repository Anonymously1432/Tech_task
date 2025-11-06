-- name: CreateWalletOperation :one
INSERT INTO requests_history (operationType, amount, wallet_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetWalletAmount :one
SELECT amount
FROM wallets
WHERE id = $1;

-- name: AddAmount :exec
UPDATE wallets
SET amount = amount + $2
WHERE id = $1;

-- name: SubtractAmountIfEnough :exec
UPDATE wallets
SET amount = amount - $2
WHERE id = $1 AND amount >= $2;