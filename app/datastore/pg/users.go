package pg

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"gamma/app/datastore"
	"gamma/app/domain/bo"
	"gamma/app/domain/definition"

	"github.com/labstack/gommon/log"
)

// TODO(add full user profile to user data)

var (
	userOnce     sync.Once
	userInstance definition.UserRepository
)

type userRepo struct {
	dbInstance *sql.DB
}

func GetUserRepo() definition.UserRepository {
	userOnce.Do(func() {
		userInstance = &userRepo{
			dbInstance: datastore.RwInstance(),
		}
	})
	return userInstance
}

func (u *userRepo) GetUser(ctx context.Context, uuid string) (*bo.User, error) {
	statement := "SELECT uuid, email, first_name, last_name, username, phone_number, image_url FROM users WHERE uuid = $1"
	res := u.dbInstance.QueryRowContext(ctx, statement, uuid)

	if res.Err() != nil {
		log.Errorf("could not get user by uuid: %s", res.Err().Error())
		return nil, res.Err()
	}

	var user bo.User

	if err := res.Scan(
		&user.Uuid,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.UserName,
		&user.PhoneNumber,
		&user.ImageUrl,
	); err != nil {
		log.Errorf("could not scan res into user object: %v", err)
		return nil, err
	}

	return &user, nil
}

func (u *userRepo) GetUserByEmail(ctx context.Context, email string) (*bo.User, error) {
	statement := "SELECT id, uuid, email, first_name, last_name, password_hash, username, phone_number, image_url FROM users WHERE email = $1"
	res := u.dbInstance.QueryRowContext(ctx, statement, email)

	if res.Err() != nil {
		log.Errorf("could not get user by uuid: %s", res.Err().Error())
		return nil, res.Err()
	}

	var user bo.User

	if err := res.Scan(
		&user.Id,
		&user.Uuid,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.PasswordHash,
		&user.UserName,
		&user.PhoneNumber,
		&user.ImageUrl,
	); err != nil {
		log.Errorf("could not scan res into user object: %s", err.Error())
		return nil, err
	}

	return &user, nil
}

func (u *userRepo) GetUserOrganizations(ctx context.Context, userId int) ([]bo.OrganizationUserJoin, error) {
	statement := "SELECT org_name, org_image_url, city, uuid, policies_num FROM org_users INNER JOIN organizations ON org_users.organization_fk = organizations.id WHERE org_users.user_fk = $1"
	res, err := u.dbInstance.QueryContext(ctx, statement, userId)

	if err != nil {
		log.Errorf("could not get org by user id: %s", err.Error())
		return nil, err
	}

	var orgUserJoins []bo.OrganizationUserJoin

	for res.Next() {
		var organizationUserJoin bo.OrganizationUserJoin

		if err := res.Scan(
			&organizationUserJoin.OrganizationName,
			&organizationUserJoin.ImageUrl,
			&organizationUserJoin.City,
			&organizationUserJoin.Uuid,
			&organizationUserJoin.PolicyNum,
		); err != nil {
			log.Errorf("error while scanning: %s", err.Error())
			return nil, err
		}

		orgUserJoins = append(orgUserJoins, organizationUserJoin)
	}

	return orgUserJoins, nil
}

func (u *userRepo) GetOrganizationEvents(ctx context.Context, orgUuid string) ([]bo.Event, error) {
	statement := "SELECT events.id, event_name, event_date, event_location, event_image_url, uuid FROM events INNER JOIN organizations ON events.organization_fk = organizations.id WHERE organizations.uuid = $1"
	res, err := u.dbInstance.QueryContext(ctx, statement, orgUuid)

	if err != nil {
		log.Errorf("could not get events by org uuid: %s", err.Error())
		return nil, err
	}

	var events []bo.Event

	for res.Next() {
		var event bo.Event

		if err := res.Scan(
			&event.Id,
			&event.EventName,
			&event.EventDate,
			&event.EventLocation,
			&event.EventImage,
			&event.Uuid,
		); err != nil {
			log.Errorf("error while scanning: %s", err.Error())
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

func (u *userRepo) InsertEventByOrganization(ctx context.Context, orgUuid string, event *bo.Event) (*bo.Event, error) {
	statement := `
		WITH org_id AS (
			SELECT id FROM organizations WHERE uuid = $1
		)
		INSERT INTO events (event_name, event_date, event_location, event_image_url, uuid, organization_fk) VALUES ($2,$3,$4,$5,$6, (SELECT id FROM org_id))
		RETURNING id
	`
	res := u.dbInstance.QueryRowContext(ctx, statement, orgUuid, event.EventName, event.EventDate, event.EventLocation, event.EventImage, event.Uuid)

	if res.Err() != nil {
		log.Errorf("could not insert event by organization uuid: %s", res.Err().Error())
		return nil, res.Err()
	}

	var id int
	res.Scan(&id)
	event.Id = id

	return event, nil
}

func (u *userRepo) UpdateEventById(ctx context.Context, eventId int, updatedEvent *bo.Event) error {
	statement := `
		UPDATE events
		SET
			event_name = $2,
			event_date = $3,
			event_location = $4,
			uuid = $5
		WHERE id = $1
	`

	_, err := u.dbInstance.ExecContext(ctx, statement, eventId, updatedEvent.EventName, updatedEvent.EventDate, updatedEvent.EventLocation, updatedEvent.Uuid)

	if err != nil {
		log.Errorf("could not update event with id %d: %s", eventId, err.Error())
	}

	return err
}

func (u *userRepo) DeleteEventById(ctx context.Context, eventId int) error {
	statement := `DELETE FROM events WHERE id = $1`
	_, err := u.dbInstance.ExecContext(ctx, statement, eventId)

	if err != nil {
		log.Errorf("could not delete event by id: %s", err.Error())
	}

	return err
}

func (u *userRepo) InsertUser(ctx context.Context, uuid, email, phone_number, hash, firstName, lastName, userName, image_url string) error {
	statement := "INSERT INTO users (uuid, email, password_hash, first_name, last_name, username, phone_number, image_url) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
	_, err := u.dbInstance.ExecContext(ctx, statement, uuid, email, hash, firstName, lastName, userName, phone_number, image_url)

	if err != nil {
		log.Errorf("could not insert user into db: %s", err.Error())
		return err
	}

	return nil
}

func (u *userRepo) GetUserEvents(ctx context.Context, userId int) ([]bo.Event, error) {
	statement := "SELECT event_name, event_date, event_location FROM user_events INNER JOIN events ON user_events.event_fk = events.id WHERE user_events.user_fk = $1"
	rows, err := u.dbInstance.QueryContext(ctx, statement, userId)

	if err != nil {
		log.Errorf("could not get user events: %s", err.Error())
		return nil, err
	}

	var events []bo.Event

	for rows.Next() {
		var event bo.Event

		if err := rows.Scan(
			&event.EventName,
			&event.EventDate,
			&event.EventLocation,
		); err != nil {
			log.Errorf("could not scan events into bo.Event: %s", err.Error())
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

func (u *userRepo) InsertUserEvent(ctx context.Context, userId, eventId int) error {
	statement := "INSERT INTO user_events (user_fk, event_fk) VALUES ($1, $2)"
	_, err := u.dbInstance.ExecContext(ctx, statement, userId, eventId)

	if err != nil {
		log.Errorf("could not insert into user_events: %s", err.Error())
		return err
	}

	return nil
}

func (u *userRepo) GetEvent(ctx context.Context, eventUuid string) (*bo.Event, error) {
	statement := "SELECT id, event_name, event_date, event_location, uuid FROM events WHERE uuid = $1"
	row := u.dbInstance.QueryRowContext(ctx, statement, eventUuid)

	if row.Err() != nil {
		log.Errorf("could not get event by uuid: %s", row.Err().Error())
		return nil, row.Err()
	}

	var event bo.Event

	if err := row.Scan(
		&event.Id,
		&event.EventName,
		&event.EventDate,
		&event.EventLocation,
		&event.Uuid,
	); err != nil {
		log.Errorf("could not scan event into bo.Event: %s", err)
		return nil, err
	}

	return &event, nil
}

func (u *userRepo) InsertOrganization(ctx context.Context, uuid, name, city, image_url string) (int, error) {
	statement := "INSERT INTO organizations (org_name, city, uuid, org_image_url) VALUES ($1, $2, $3, $4) RETURNING id"
	res := u.dbInstance.QueryRowContext(ctx, statement, name, city, uuid, image_url)

	if res.Err() != nil {
		log.Errorf("could not insert org: %s", res.Err().Error())
		return -1, res.Err()
	}

	var id int
	if err := res.Scan(&id); err != nil {
		log.Errorf("could not scan id: %s", err.Error())
		return -1, err
	}

	return id, nil
}

func (u *userRepo) InsertOrgUser(ctx context.Context, userId int, organizationId int, policiesNum int) error {
	statement := `
		INSERT INTO org_users (user_fk, organization_fk, policies_num) VALUES ($1,$2,$3)
	`
	_, err := u.dbInstance.ExecContext(ctx, statement, userId, organizationId, policiesNum)

	if err != nil {
		log.Errorf("could not insert org user: %s", err.Error())
		return err
	}

	return nil
}

func (u *userRepo) InsertEvent(ctx context.Context, event_name string, event_date time.Time, event_location string, uuid string, organizationFk int) error {
	statement := "INSERT INTO events (event_name, event_date, event_location, uuid, organization_fk) VALUES ($1, $2, $3, $4, $5)"
	_, err := u.dbInstance.ExecContext(ctx, statement, event_name, event_date, event_location, uuid, organizationFk)

	if err != nil {
		log.Errorf("could not insert into events: %s", err.Error())
		return err
	}

	return nil
}

func (u *userRepo) InsertInvite(ctx context.Context, expirationDate time.Time, useLimit int, uuid string, policy bo.InvitePolicy) error {
	statement := "INSERT INTO invites (expiration_date, use_limit, uuid, policy_json) VALUES ($1, $2, $3, $4)"
	_, err := u.dbInstance.ExecContext(ctx, statement, expirationDate, useLimit, uuid, policy)

	if err != nil {
		log.Errorf("could not insert into invites: %s", err.Error())
		return err
	}

	return nil
}

func (u *userRepo) GetInvite(ctx context.Context, inviteUuid string) (*bo.Invite, error){
	statement := "SELECT id, expiration_date, use_limit, policy_json FROM invites WHERE uuid = $1"
	row := u.dbInstance.QueryRowContext(ctx, statement, inviteUuid)

	if row.Err() != nil {
		log.Errorf("could not get invite by uuid: %s", row.Err().Error())
		return nil, row.Err()
	}

	var invite bo.Invite

	if err := row.Scan(
		&invite.Id,
		&invite.ExpirationDate,
		&invite.UseLimit,
		&invite.Policy,
	); err != nil {
		log.Errorf("could not scan event into bo.Invite: %s", err.Error())
		return nil, err
	}

	return &invite, nil
}

func (u *userRepo) DecrementInvite(ctx context.Context, inviteId int) error {
	statement := `
	UPDATE invites
	SET 
		use_limit = use_limit - 1
	WHERE id = $1
	`

	_, err := u.dbInstance.ExecContext(ctx, statement, inviteId)
	if err != nil {
		log.Errorf("could not decrement invite use limit: %s", err.Error())
		return err
	}

	return nil
}
