package pg

import (
	"gamma/app/datastore"

	"github.com/labstack/gommon/log"
)

func ClearAll() {
	_, err := datastore.RwInstance().Exec("TRUNCATE users, org_users, organizations, events, user_events, user_event_invites CASCADE")
	if err != nil {
		log.Errorf("could not wipe database clean: %s", err.Error())
	}
}
