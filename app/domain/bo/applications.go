package bo

import "database/sql"

type ApplicationState string

const (
	APPLIED  ApplicationState = "APPLIED"
	ACCEPTED ApplicationState = "ACCEPTED"
	NULL     ApplicationState = "NULL"
)

func ParseApplicationState(state sql.NullString) ApplicationState {
	if state.Valid {
		switch state.String {
		case "APPLIED":
			return APPLIED
		case "ACCEPTED":
			return ACCEPTED
		}
	}
	return NULL
}
