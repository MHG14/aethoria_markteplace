-- name: GetGuild :one
SELECT * FROM guilds WHERE id = $1;

-- name: GetGuildForUpdate :one
SELECT * FROM guilds WHERE id = $1 FOR UPDATE;

-- name: CreateGuild :one
INSERT INTO guilds (name, total_money, daily_limit)
VALUES ($1, $2, $3)
    RETURNING *;

-- name: UpdateGuildWallet :one
UPDATE guilds
SET total_money    = $2,
    reserved_money = $3,
    daily_spent    = $4
WHERE id = $1
    RETURNING *;

-- name: TopUpGuildWallet :one
UPDATE guilds
SET total_money = total_money + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
    RETURNING *;

-- name: ResetDailySpent :exec
UPDATE guilds SET daily_spent = 0;