-- name: CreateDelegation :one
INSERT INTO delegations (delegator, timestamp, amount, level)
VALUES ($1, $2, $3, $4)
RETURNING delegator, timestamp, amount, level;
