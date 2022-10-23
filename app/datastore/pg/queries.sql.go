// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: queries.sql

package userRepo

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const batchAddOrgUsersToEvent = `-- name: BatchAddOrgUsersToEvent :exec
INSERT INTO user_events (user_fk, event_fk, application_state) VALUES (
    unnest($1::uuid[]),
    unnest($2::uuid[]),
    1
)
`

type BatchAddOrgUsersToEventParams struct {
	UserUuids  []uuid.UUID
	EventUuids []uuid.UUID
}

func (q *Queries) BatchAddOrgUsersToEvent(ctx context.Context, arg *BatchAddOrgUsersToEventParams) error {
	_, err := q.db.Exec(ctx, batchAddOrgUsersToEvent, arg.UserUuids, arg.EventUuids)
	return err
}

const checkUser = `-- name: CheckUser :one
SELECT id, user_fk, event_fk, application_state FROM user_events ue WHERE ue.user_fk = $1 AND ue.event_fk = $2 AND ue.application_state = 2
`

type CheckUserParams struct {
	UserFk  uuid.UUID
	EventFk uuid.UUID
}

func (q *Queries) CheckUser(ctx context.Context, arg *CheckUserParams) (*UserEvent, error) {
	row := q.db.QueryRow(ctx, checkUser, arg.UserFk, arg.EventFk)
	var i UserEvent
	err := row.Scan(
		&i.ID,
		&i.UserFk,
		&i.EventFk,
		&i.ApplicationState,
	)
	return &i, err
}

const getEvent = `-- name: GetEvent :one
SELECT id, event_name, event_date, event_location, event_description, event_image_url, organization_fk FROM events e WHERE id = $1
`

func (q *Queries) GetEvent(ctx context.Context, id uuid.UUID) (*Event, error) {
	row := q.db.QueryRow(ctx, getEvent, id)
	var i Event
	err := row.Scan(
		&i.ID,
		&i.EventName,
		&i.EventDate,
		&i.EventLocation,
		&i.EventDescription,
		&i.EventImageUrl,
		&i.OrganizationFk,
	)
	return &i, err
}

const getEventsWithOrganizations = `-- name: GetEventsWithOrganizations :many
SELECT e.id, event_name, event_date, event_location, event_description, event_image_url, organization_fk, o.id, org_name, city, org_image_url FROM events e INNER JOIN organizations o ON e.organization_fk = o.id
WHERE e.event_date > $1
    AND (CASE WHEN $2::bool THEN o.id = $3 ELSE TRUE END)
ORDER BY e.event_date DESC
`

type GetEventsWithOrganizationsParams struct {
	DateFloor          time.Time
	FilterOrganization bool
	OrgUuid            uuid.UUID
}

type GetEventsWithOrganizationsRow struct {
	ID               uuid.UUID
	EventName        string
	EventDate        time.Time
	EventLocation    string
	EventDescription string
	EventImageUrl    string
	OrganizationFk   uuid.UUID
	ID_2             uuid.UUID
	OrgName          string
	City             string
	OrgImageUrl      string
}

func (q *Queries) GetEventsWithOrganizations(ctx context.Context, arg *GetEventsWithOrganizationsParams) ([]*GetEventsWithOrganizationsRow, error) {
	rows, err := q.db.Query(ctx, getEventsWithOrganizations, arg.DateFloor, arg.FilterOrganization, arg.OrgUuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetEventsWithOrganizationsRow
	for rows.Next() {
		var i GetEventsWithOrganizationsRow
		if err := rows.Scan(
			&i.ID,
			&i.EventName,
			&i.EventDate,
			&i.EventLocation,
			&i.EventDescription,
			&i.EventImageUrl,
			&i.OrganizationFk,
			&i.ID_2,
			&i.OrgName,
			&i.City,
			&i.OrgImageUrl,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getInvite = `-- name: GetInvite :one
SELECT id, expiration_date, capacity, user_fk, org_fk, entity_uuid, entity_type FROM invites i WHERE id = $1
`

func (q *Queries) GetInvite(ctx context.Context, id uuid.UUID) (*Invite, error) {
	row := q.db.QueryRow(ctx, getInvite, id)
	var i Invite
	err := row.Scan(
		&i.ID,
		&i.ExpirationDate,
		&i.Capacity,
		&i.UserFk,
		&i.OrgFk,
		&i.EntityUuid,
		&i.EntityType,
	)
	return &i, err
}

const getInvitesForOrgUser = `-- name: GetInvitesForOrgUser :many
SELECT id, expiration_date, capacity, user_fk, org_fk, entity_uuid, entity_type FROM invites i WHERE user_fk = $1
`

func (q *Queries) GetInvitesForOrgUser(ctx context.Context, userFk uuid.UUID) ([]*Invite, error) {
	rows, err := q.db.Query(ctx, getInvitesForOrgUser, userFk)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*Invite
	for rows.Next() {
		var i Invite
		if err := rows.Scan(
			&i.ID,
			&i.ExpirationDate,
			&i.Capacity,
			&i.UserFk,
			&i.OrgFk,
			&i.EntityUuid,
			&i.EntityType,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getOrgUser = `-- name: GetOrgUser :one
SELECT id, policies_num, user_fk, organization_fk FROM org_users ou WHERE ou.user_fk = $1 AND ou.organization_fk = $2 LIMIT 1
`

type GetOrgUserParams struct {
	UserFk         uuid.UUID
	OrganizationFk uuid.UUID
}

func (q *Queries) GetOrgUser(ctx context.Context, arg *GetOrgUserParams) (*OrgUser, error) {
	row := q.db.QueryRow(ctx, getOrgUser, arg.UserFk, arg.OrganizationFk)
	var i OrgUser
	err := row.Scan(
		&i.ID,
		&i.PoliciesNum,
		&i.UserFk,
		&i.OrganizationFk,
	)
	return &i, err
}

const getOrganization = `-- name: GetOrganization :one
SELECT id, org_name, city, org_image_url FROM organizations o WHERE id = $1
`

func (q *Queries) GetOrganization(ctx context.Context, id uuid.UUID) (*Organization, error) {
	row := q.db.QueryRow(ctx, getOrganization, id)
	var i Organization
	err := row.Scan(
		&i.ID,
		&i.OrgName,
		&i.City,
		&i.OrgImageUrl,
	)
	return &i, err
}

const getOrganizationByUuid = `-- name: GetOrganizationByUuid :one
SELECT id, org_name, city, org_image_url FROM organizations o WHERE o.id = $1 LIMIT 1
`

func (q *Queries) GetOrganizationByUuid(ctx context.Context, organizationUuid uuid.UUID) (*Organization, error) {
	row := q.db.QueryRow(ctx, getOrganizationByUuid, organizationUuid)
	var i Organization
	err := row.Scan(
		&i.ID,
		&i.OrgName,
		&i.City,
		&i.OrgImageUrl,
	)
	return &i, err
}

const getOrganizationEvents = `-- name: GetOrganizationEvents :many
SELECT id, event_name, event_date, event_location, event_description, event_image_url, organization_fk FROM events e WHERE e.organization_fk = $1
`

func (q *Queries) GetOrganizationEvents(ctx context.Context, orgUuid uuid.UUID) ([]*Event, error) {
	rows, err := q.db.Query(ctx, getOrganizationEvents, orgUuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*Event
	for rows.Next() {
		var i Event
		if err := rows.Scan(
			&i.ID,
			&i.EventName,
			&i.EventDate,
			&i.EventLocation,
			&i.EventDescription,
			&i.EventImageUrl,
			&i.OrganizationFk,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getOrganizationUsers = `-- name: GetOrganizationUsers :one
SELECT ou.id, policies_num, user_fk, organization_fk, u.id, email, password_hash, phone_number, first_name, last_name, username, image_url, validated FROM org_users OU INNER JOIN users u ON ou.user_fk = u.id WHERE ou.organization_fk = $1
`

type GetOrganizationUsersRow struct {
	ID             int32
	PoliciesNum    int32
	UserFk         uuid.UUID
	OrganizationFk uuid.UUID
	ID_2           uuid.UUID
	Email          string
	PasswordHash   string
	PhoneNumber    string
	FirstName      string
	LastName       string
	Username       string
	ImageUrl       string
	Validated      bool
}

func (q *Queries) GetOrganizationUsers(ctx context.Context, orgUuid uuid.UUID) (*GetOrganizationUsersRow, error) {
	row := q.db.QueryRow(ctx, getOrganizationUsers, orgUuid)
	var i GetOrganizationUsersRow
	err := row.Scan(
		&i.ID,
		&i.PoliciesNum,
		&i.UserFk,
		&i.OrganizationFk,
		&i.ID_2,
		&i.Email,
		&i.PasswordHash,
		&i.PhoneNumber,
		&i.FirstName,
		&i.LastName,
		&i.Username,
		&i.ImageUrl,
		&i.Validated,
	)
	return &i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, email, password_hash, phone_number, first_name, last_name, username, image_url, validated FROM users WHERE email = $1 LIMIT 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.PasswordHash,
		&i.PhoneNumber,
		&i.FirstName,
		&i.LastName,
		&i.Username,
		&i.ImageUrl,
		&i.Validated,
	)
	return &i, err
}

const getUserByUuid = `-- name: GetUserByUuid :one

SELECT id, email, password_hash, phone_number, first_name, last_name, username, image_url, validated FROM users WHERE id = $1 LIMIT 1
`

// GETS
func (q *Queries) GetUserByUuid(ctx context.Context, uuid uuid.UUID) (*User, error) {
	row := q.db.QueryRow(ctx, getUserByUuid, uuid)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.PasswordHash,
		&i.PhoneNumber,
		&i.FirstName,
		&i.LastName,
		&i.Username,
		&i.ImageUrl,
		&i.Validated,
	)
	return &i, err
}

const getUserEvents = `-- name: GetUserEvents :many
SELECT ue.id, user_fk, event_fk, application_state, e.id, event_name, event_date, event_location, event_description, event_image_url, organization_fk FROM user_events ue INNER JOIN events e ON ue.event_fk = e.id WHERE ue.user_fk = $1
`

type GetUserEventsRow struct {
	ID               int32
	UserFk           uuid.UUID
	EventFk          uuid.UUID
	ApplicationState int32
	ID_2             uuid.UUID
	EventName        string
	EventDate        time.Time
	EventLocation    string
	EventDescription string
	EventImageUrl    string
	OrganizationFk   uuid.UUID
}

func (q *Queries) GetUserEvents(ctx context.Context, userUuid uuid.UUID) ([]*GetUserEventsRow, error) {
	rows, err := q.db.Query(ctx, getUserEvents, userUuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetUserEventsRow
	for rows.Next() {
		var i GetUserEventsRow
		if err := rows.Scan(
			&i.ID,
			&i.UserFk,
			&i.EventFk,
			&i.ApplicationState,
			&i.ID_2,
			&i.EventName,
			&i.EventDate,
			&i.EventLocation,
			&i.EventDescription,
			&i.EventImageUrl,
			&i.OrganizationFk,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserOrganizations = `-- name: GetUserOrganizations :many
SELECT ou.id, policies_num, user_fk, organization_fk, og.id, org_name, city, org_image_url FROM org_users ou INNER JOIN organizations og ON ou.organization_fk = og.id WHERE ou.user_fk = $1
`

type GetUserOrganizationsRow struct {
	ID             int32
	PoliciesNum    int32
	UserFk         uuid.UUID
	OrganizationFk uuid.UUID
	ID_2           uuid.UUID
	OrgName        string
	City           string
	OrgImageUrl    string
}

func (q *Queries) GetUserOrganizations(ctx context.Context, userUuid uuid.UUID) ([]*GetUserOrganizationsRow, error) {
	rows, err := q.db.Query(ctx, getUserOrganizations, userUuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetUserOrganizationsRow
	for rows.Next() {
		var i GetUserOrganizationsRow
		if err := rows.Scan(
			&i.ID,
			&i.PoliciesNum,
			&i.UserFk,
			&i.OrganizationFk,
			&i.ID_2,
			&i.OrgName,
			&i.City,
			&i.OrgImageUrl,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserWithOrg = `-- name: GetUserWithOrg :one
SELECT ou.id, policies_num, user_fk, organization_fk, u.id, email, password_hash, phone_number, first_name, last_name, username, image_url, validated FROM org_users ou INNER JOIN  users u ON ou.user_fk = u.id WHERE ou.organization_fk = $1 AND u.id = $2 LIMIT 1
`

type GetUserWithOrgParams struct {
	OrgUuid  uuid.UUID
	UserUuid uuid.UUID
}

type GetUserWithOrgRow struct {
	ID             int32
	PoliciesNum    int32
	UserFk         uuid.UUID
	OrganizationFk uuid.UUID
	ID_2           uuid.UUID
	Email          string
	PasswordHash   string
	PhoneNumber    string
	FirstName      string
	LastName       string
	Username       string
	ImageUrl       string
	Validated      bool
}

func (q *Queries) GetUserWithOrg(ctx context.Context, arg *GetUserWithOrgParams) (*GetUserWithOrgRow, error) {
	row := q.db.QueryRow(ctx, getUserWithOrg, arg.OrgUuid, arg.UserUuid)
	var i GetUserWithOrgRow
	err := row.Scan(
		&i.ID,
		&i.PoliciesNum,
		&i.UserFk,
		&i.OrganizationFk,
		&i.ID_2,
		&i.Email,
		&i.PasswordHash,
		&i.PhoneNumber,
		&i.FirstName,
		&i.LastName,
		&i.Username,
		&i.ImageUrl,
		&i.Validated,
	)
	return &i, err
}

const insertEvent = `-- name: InsertEvent :exec
INSERT INTO events (id, event_name, event_date, event_location, event_description, event_image_url, organization_fk) VALUES ($1,$2,$3,$4,$5,$6,$7)
`

type InsertEventParams struct {
	ID               uuid.UUID
	EventName        string
	EventDate        time.Time
	EventLocation    string
	EventDescription string
	EventImageUrl    string
	OrganizationFk   uuid.UUID
}

func (q *Queries) InsertEvent(ctx context.Context, arg *InsertEventParams) error {
	_, err := q.db.Exec(ctx, insertEvent,
		arg.ID,
		arg.EventName,
		arg.EventDate,
		arg.EventLocation,
		arg.EventDescription,
		arg.EventImageUrl,
		arg.OrganizationFk,
	)
	return err
}

const insertInvite = `-- name: InsertInvite :exec
INSERT INTO invites (id, expiration_date, capacity, user_fk, org_fk, entity_uuid, entity_type) VALUES ($1,$2,$3,$4,$5,$6,$7)
`

type InsertInviteParams struct {
	ID             uuid.UUID
	ExpirationDate time.Time
	Capacity       int32
	UserFk         uuid.UUID
	OrgFk          uuid.UUID
	EntityUuid     uuid.UUID
	EntityType     int32
}

func (q *Queries) InsertInvite(ctx context.Context, arg *InsertInviteParams) error {
	_, err := q.db.Exec(ctx, insertInvite,
		arg.ID,
		arg.ExpirationDate,
		arg.Capacity,
		arg.UserFk,
		arg.OrgFk,
		arg.EntityUuid,
		arg.EntityType,
	)
	return err
}

const insertOrgUser = `-- name: InsertOrgUser :exec
INSERT INTO org_users (policies_num, user_fk, organization_fk) VALUES ($1,$2,$3)
`

type InsertOrgUserParams struct {
	PoliciesNum    int32
	UserFk         uuid.UUID
	OrganizationFk uuid.UUID
}

func (q *Queries) InsertOrgUser(ctx context.Context, arg *InsertOrgUserParams) error {
	_, err := q.db.Exec(ctx, insertOrgUser, arg.PoliciesNum, arg.UserFk, arg.OrganizationFk)
	return err
}

const insertOrganization = `-- name: InsertOrganization :exec
INSERT INTO organizations (id, org_name, city, org_image_url) VALUES ($1,$2,$3,$4)
`

type InsertOrganizationParams struct {
	ID          uuid.UUID
	OrgName     string
	City        string
	OrgImageUrl string
}

func (q *Queries) InsertOrganization(ctx context.Context, arg *InsertOrganizationParams) error {
	_, err := q.db.Exec(ctx, insertOrganization,
		arg.ID,
		arg.OrgName,
		arg.City,
		arg.OrgImageUrl,
	)
	return err
}

const insertUser = `-- name: InsertUser :exec

INSERT INTO users (id, email, password_hash, phone_number, first_name, last_name, username, image_url, validated) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
`

type InsertUserParams struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	PhoneNumber  string
	FirstName    string
	LastName     string
	Username     string
	ImageUrl     string
	Validated    bool
}

// PUTS
func (q *Queries) InsertUser(ctx context.Context, arg *InsertUserParams) error {
	_, err := q.db.Exec(ctx, insertUser,
		arg.ID,
		arg.Email,
		arg.PasswordHash,
		arg.PhoneNumber,
		arg.FirstName,
		arg.LastName,
		arg.Username,
		arg.ImageUrl,
		arg.Validated,
	)
	return err
}

const insertUserEvent = `-- name: InsertUserEvent :exec
INSERT INTO user_events (user_fk, event_fk, application_state) VALUES ($1,$2,$3)
`

type InsertUserEventParams struct {
	UserFk           uuid.UUID
	EventFk          uuid.UUID
	ApplicationState int32
}

func (q *Queries) InsertUserEvent(ctx context.Context, arg *InsertUserEventParams) error {
	_, err := q.db.Exec(ctx, insertUserEvent, arg.UserFk, arg.EventFk, arg.ApplicationState)
	return err
}

const truncateAll = `-- name: TruncateAll :exec

TRUNCATE users, org_users, organizations, events, user_events, invites
`

// UTIL
func (q *Queries) TruncateAll(ctx context.Context) error {
	_, err := q.db.Exec(ctx, truncateAll)
	return err
}

const updatePassword = `-- name: UpdatePassword :exec
UPDATE users SET password_hash = $1 WHERE id = $2
`

type UpdatePasswordParams struct {
	PasswordHash string
	ID           uuid.UUID
}

func (q *Queries) UpdatePassword(ctx context.Context, arg *UpdatePasswordParams) error {
	_, err := q.db.Exec(ctx, updatePassword, arg.PasswordHash, arg.ID)
	return err
}

const useInvite = `-- name: UseInvite :exec

UPDATE invites SET use_limit = use_limit - 1 WHERE id = $1 AND use_limit != 0
`

// UPDATES
func (q *Queries) UseInvite(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, useInvite, id)
	return err
}
