-- GETS

-- name: GetUserByUuid :one
SELECT * FROM users WHERE uuid = sqlc.arg(uuid)::text LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = sqlc.arg(email)::text LIMIT 1;

-- name: GetUserOrgUserJoin :one
SELECT * FROM users u INNER JOIN org_users o ON u.id = o.user_fk WHERE u.uuid = sqlc.arg(uuid)::text LIMIT 1;

-- name: GetUserOrganizations :many
SELECT * FROM org_users ou INNER JOIN organizations og ON ou.organization_fk = og.id WHERE ou.user_fk = sqlc.arg(user_id)::int;

-- name: GetOrganizationEvents :many
SELECT * FROM events e INNER JOIN organizations o ON e.organization_fk = o.id WHERE o.uuid = sqlc.arg(org_uuid)::text;

-- name: GetUserEvents :many
SELECT * FROM user_events ue INNER JOIN events e ON ue.event_fk = e.id WHERE ue.user_fk = sqlc.arg(user_id)::int;

-- name: GetEventByUuid :one
SELECT * FROM events WHERE uuid = sqlc.arg(event_uuid)::text LIMIT 1;

-- name: GetEventById :one
SELECT * FROM events WHERE id = sqlc.arg(event_id)::int LIMIT 1;

-- name: GetOrganizationByUuid :one
SELECT * FROM organizations WHERE uuid = sqlc.arg(organization_uuid)::text LIMIT 1;

-- PUTS

-- name: InsertUser :exec
INSERT INTO users (uuid, email, password_hash, phone_number, first_name, last_name, image_url, validated, refresh_token) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9);

-- name: InsertOrganization :exec
INSERT INTO organizations (org_name, city, uuid, org_image_url) VALUES ($1,$2,$3,$4);

-- name: InsertOrgUser :exec
INSERT INTO org_users (policies_num, user_fk, organization_fk) VALUES ($1,$2,$3);

-- name: InsertEvent :exec
INSERT INTO events (event_name, event_date, event_location, event_description, uuid, event_image_url, organization_fk) VALUES ($1,$2,$3,$4,$5,$6,$7);

-- UTIL

-- name: TruncateAll :exec
TRUNCATE users, org_users, organizations, events, user_events, event_applications, invites;