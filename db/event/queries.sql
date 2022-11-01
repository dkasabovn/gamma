-- GETS

-- name: GetUserByUuid :one
SELECT * FROM users WHERE id = sqlc.arg(uuid) LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = sqlc.arg(email) LIMIT 1;

-- name: GetOrgUser :one
SELECT * FROM org_users ou WHERE ou.user_fk = $1 AND ou.organization_fk = $2 LIMIT 1;

-- name: GetUserWithOrg :one
SELECT * FROM org_users ou INNER JOIN  users u ON ou.user_fk = u.id WHERE ou.organization_fk = sqlc.arg(org_uuid) AND u.id = sqlc.arg(user_uuid) LIMIT 1;

-- name: GetUserOrganizations :many
SELECT * FROM org_users ou INNER JOIN organizations og ON ou.organization_fk = og.id WHERE ou.user_fk = sqlc.arg(user_uuid);

-- name: GetOrganizationEvents :many
SELECT * FROM events e WHERE e.organization_fk = sqlc.arg(org_uuid);

-- name: GetUserEvents :many
SELECT * FROM user_events ue INNER JOIN events e ON ue.event_fk = e.id WHERE ue.user_fk = sqlc.arg(user_uuid);

-- name: GetOrganizationByUuid :one
SELECT * FROM organizations o WHERE o.id = sqlc.arg(organization_uuid) LIMIT 1;

-- name: CheckUser :one
SELECT * FROM user_events ue WHERE ue.user_fk = $1 AND ue.event_fk = $2 AND ue.application_state = 2;

-- name: GetEventsWithOrganizations :many
SELECT * FROM events e INNER JOIN organizations o ON e.organization_fk = o.id
WHERE e.event_date > @date_floor
    AND (CASE WHEN @filter_organization::bool THEN o.id = @org_uuid ELSE TRUE END)
ORDER BY e.event_date DESC;

-- name: GetInvitesForOrgUser :many
SELECT * FROM invites i WHERE user_fk = $1;

-- name: GetInvite :one
SELECT * FROM invites i WHERE id = $1;

-- name: GetEvent :one
SELECT * FROM events e WHERE id = $1;

-- name: GetOrganization :one
SELECT * FROM organizations o WHERE id = $1;

-- name: GetOrganizationUsers :one
SELECT * FROM org_users OU INNER JOIN users u ON ou.user_fk = u.id WHERE ou.organization_fk = sqlc.arg(org_uuid);

-- PUTS

-- name: InsertUser :exec
INSERT INTO users (id, email, password_hash, phone_number, first_name, last_name, username, image_url, validated) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9);

-- name: InsertOrganization :exec
INSERT INTO organizations (id, org_name, city, org_image_url) VALUES ($1,$2,$3,$4);

-- name: InsertOrgUser :exec
INSERT INTO org_users (policies_num, user_fk, organization_fk) VALUES ($1,$2,$3);

-- name: InsertEvent :exec
INSERT INTO events (id, event_name, event_date, event_location, event_description, event_image_url, organization_fk) VALUES ($1,$2,$3,$4,$5,$6,$7);

-- name: InsertInvite :exec
INSERT INTO invites (id, expiration_date, capacity, user_fk, org_fk, entity_uuid, entity_type) VALUES ($1,$2,$3,$4,$5,$6,$7);

-- name: InsertUserEvent :exec
INSERT INTO user_events (user_fk, event_fk, application_state) VALUES ($1,$2,$3);

-- name: BatchAddOrgUsersToEvent :exec
INSERT INTO user_events (user_fk, event_fk, application_state) VALUES (
    unnest(@user_uuids::uuid[]),
    unnest(@event_uuids::uuid[]),
    1
);

-- UPDATES

-- name: UseInvite :exec
UPDATE invites SET use_limit = use_limit - 1 WHERE id = $1 AND use_limit != 0;

-- name: UpdatePassword :exec
UPDATE users SET password_hash = $1 WHERE id = $2;

-- name: UpdateUser :exec
UPDATE users SET first_name = $1, last_name = $2, email = $3, phone_number = $4, username = $5, image_url = $6 WHERE id = $7;

-- UTIL

-- name: TruncateAll :exec
TRUNCATE users, org_users, organizations, events, user_events, invites;
