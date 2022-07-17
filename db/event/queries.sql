-- GETS

-- name: GetUserByUuid :one
SELECT * FROM users WHERE uuid = sqlc.arg(uuid)::text LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = sqlc.arg(email)::text LIMIT 1;

-- name: GetOrgUser :one
SELECT * FROM users u INNER JOIN org_users o ON u.id = o.user_fk INNER JOIN organizations org ON org.id = o.organization_fk WHERE u.uuid = sqlc.arg(user_uuid)::text AND org.uuid = sqlc.arg(org_uuid)::text LIMIT 1;

-- name: GetUserOrganizations :many
SELECT * FROM org_users ou INNER JOIN organizations og ON ou.organization_fk = og.id WHERE ou.user_fk = sqlc.arg(user_id)::int;

-- name: GetOrganizationEvents :many
SELECT e.id, e.event_name, e.event_date, e.event_location, e.event_description, e.uuid, e.event_image_url, e.organization_fk FROM events e INNER JOIN organizations o ON e.organization_fk = o.id WHERE o.uuid = sqlc.arg(org_uuid)::text;

-- name: GetUserEvents :many
SELECT * FROM user_events ue INNER JOIN events e ON ue.event_fk = e.id WHERE ue.user_fk = sqlc.arg(user_id)::int;

-- name: GetOrganizationByUuid :one
SELECT * FROM organizations WHERE uuid = sqlc.arg(organization_uuid)::text LIMIT 1;

-- name: GetEvents :many
SELECT * FROM events e INNER JOIN organizations o ON o.id = e.organization_fk WHERE event_date > NOW() ORDER BY event_date - NOW() ASC LIMIT 50;

-- name: SearchEvents :many
SELECT * FROM events e INNER JOIN organizations o ON o.id = e.organization_fk WHERE event_name LIKE sqlc.arg(event_name_like_query)::text LIMIT 10;

-- PUTS

-- name: InsertUser :exec
INSERT INTO users (uuid, email, password_hash, phone_number, first_name, last_name, image_url, validated, refresh_token) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9);

-- name: InsertOrganization :one
INSERT INTO organizations (org_name, city, uuid, org_image_url) VALUES ($1,$2,$3,$4) RETURNING id;

-- name: InsertOrgUser :exec
INSERT INTO org_users (policies_num, user_fk, organization_fk) VALUES ($1,$2,$3);

-- name: InsertEvent :exec
INSERT INTO events (event_name, event_date, event_location, event_description, uuid, event_image_url, organization_fk) VALUES ($1,$2,$3,$4,$5,$6,$7);

-- name: InsertInvite :exec
INSERT INTO invites (expiration_date, capacity, policy_json, uuid, org_user_uuid, org_fk) VALUES ($1, $2, $3, $4, $5, $6);

-- UPDATES

UPDATE invites SET use_limit = use_limit - 1 WHERE id = $1 AND use_limit > 0;

-- UTIL

-- name: TruncateAll :exec
TRUNCATE users, org_users, organizations, events, user_events, event_applications, invites;