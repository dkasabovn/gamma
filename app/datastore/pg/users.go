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
	statement := "SELECT id, uuid, email, first_name, last_name, username, org_user_fk FROM users LEFT JOIN org_users ON org_user_fk = org_users.id WHERE uuid = $1"
	res := u.dbInstance.QueryRowContext(ctx, statement, uuid)

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
		&user.UserName,
		&user.OrgUserFk,
	); err != nil {
		log.Errorf("could not scan res into user object: %v", err)
		return nil, err
	}

	return &user, nil
}

func (u *userRepo) GetUserByEmail(ctx context.Context, email string) (*bo.User, error) {
	statement := "SELECT id, uuid, email, first_name, last_name, password_hash, org_user_fk FROM users WHERE email = $1"
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
		&user.OrgUserFk,
	); err != nil {
		log.Errorf("could not scan res into user object: %s", err.Error())
		return nil, err
	}

	return &user, nil
}

func (u *userRepo) InsertUser(ctx context.Context, uuid string, email string, hash string, firstName string, lastName string, userName string) error {
	statement := "INSERT INTO users (uuid, email, password_hash, first_name, last_name, username) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := u.dbInstance.ExecContext(ctx, statement, uuid, email, hash, firstName, lastName, userName)

	if err != nil {
		log.Errorf("could not insert user into db: %s", err.Error())
		return err
	}

	return nil
}

func (u *userRepo) GetOrgUserEvents(ctx context.Context, orgUserId sql.NullInt64) ([]bo.Event, error) {
	statement := "SELECT event_name, event_date, event_location, uuid FROM org_users INNER JOIN events ON org_users.organization_fk = events.organization_fk WHERE org_users.id = $1"
	rows, err := u.dbInstance.QueryContext(ctx, statement, orgUserId)

	if err != nil {
		log.Errorf("could not get events from org user: %s", err.Error())
		return nil, err
	}

	var events []bo.Event

	for rows.Next() {
		var event bo.Event
		if err := rows.Scan(
			&event.EventName,
			&event.EventDate,
			&event.EventLocation,
			&event.Uuid,
		); err != nil {
			log.Errorf("error scanning event into bo.Event: %s", err.Error())
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

func (u *userRepo) DeleteInvite(ctx context.Context, inviteId int) error {
	statement := "DELETE FROM user_event_invites WHERE id = $1"
	_, err := u.dbInstance.ExecContext(ctx, statement, inviteId)

	if err != nil {
		log.Errorf("could not delete invite: %s", err.Error())
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

func (u *userRepo) InsertOrganization(ctx context.Context, uuid, name, city string) (int, error) {
	statement := "INSERT INTO organizations (org_name, city, uuid) VALUES ($1, $2, $3) RETURNING id"
	res := u.dbInstance.QueryRowContext(ctx, statement, name, city, uuid)

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

func (u *userRepo) InsertOrgUser(ctx context.Context, userUuid string, organizationId int) error {
	statement := `
		WITH new_org_user AS (
			INSERT INTO
				org_users (organization_fk)
			VALUES
				($1)
			RETURNING id
		)
		UPDATE
			users
		SET
			org_user_fk = (SELECT id FROM new_org_user)
		WHERE
			uuid = $2
	`
	_, err := u.dbInstance.ExecContext(ctx, statement, organizationId, userUuid)

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
