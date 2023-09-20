-- name: GetUser :one
SELECT * FROM users WHERE id = ? LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = ? LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = ? LIMIT 1;

-- name: GetUserByConfirmationToken :one
SELECT * FROM users WHERE confirm_token = ? LIMIT 1;

-- name: GetUserByRecoveryToken :one
SELECT * FROM users WHERE recovery_token = ? LIMIT 1;

-- name: GetUserByOAuth2Identity :one
SELECT u.* FROM users u JOIN user_identities i ON u.id = i.user_id WHERE i.provider = ? AND i.identity = ? LIMIT 1;

-- name: GetUserRoles :many
SELECT role FROM user_roles WHERE user_id = ?;

-- name: GetUserPermissions :many
SELECT permission FROM user_permissions WHERE user_id = ?;

-- name: CreateUser :execresult
INSERT INTO users (email, password, totp_secret, confirmed, confirm_token, recovery_token, locked_until, attempts, last_attempt) VALUES (?,?,?,?,?,?,?,?,?);

-- name: UpdateUser :exec
UPDATE users SET email = ?, password = ?, totp_secret = ?, confirmed = ?, confirm_token = ?, recovery_token = ?, locked_until = ?, attempts = ?, last_attempt = ? WHERE id = ?;

-- name: CreateOAuth2Identity :exec
INSERT INTO user_identities (user_id, provider, identity) VALUES (?,?,?);

-- name: RemoveOAuth2Identity :exec
DELETE FROM user_identities WHERE user_id = ? AND provider = ? AND identity = ? LIMIT 1;

-- name: GetAccessToken :one
SELECT * FROM access_tokens WHERE user_id = ? AND token = ? LIMIT 1;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens WHERE user_id = ? AND token = ? LIMIT 1;

-- name: CreateAccessToken :exec
REPLACE INTO access_tokens (user_id, token, valid, chain, created_at) VALUES (?,?,?,?,?);

-- name: CreateRefreshToken :exec
REPLACE INTO refresh_tokens (user_id, token, valid, chain, created_at) VALUES (?,?,?,?,?);

-- name: InvalidateAccessTokens :exec
UPDATE access_tokens SET valid = 0 WHERE user_id = ?;

-- name: InvalidateRefreshToken :exec
UPDATE refresh_tokens SET valid = 0 WHERE user_id = ? AND token = ? LIMIT 1;

-- name: InvalidateRefreshTokens :exec
UPDATE refresh_tokens SET valid = 0 WHERE user_id = ?;

-- name: InvalidateAccessTokenChain :exec
UPDATE access_tokens SET valid = 0 WHERE user_id = ? AND chain = ?;

-- name: InvalidateRefreshTokenChain :exec
UPDATE refresh_tokens SET valid = 0 WHERE user_id = ? AND chain = ?;

-- name: PruneAccessTokens :exec
DELETE FROM access_tokens WHERE created_at < ?;

-- name: PruneRefreshTokens :exec
DELETE FROM refresh_tokens WHERE created_at < ?;