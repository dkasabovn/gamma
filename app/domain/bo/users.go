package bo

import "database/sql"

type User struct {
	Id           int
	Uuid         string
	Email        string
	FirstName    string
	LastName     string
	PasswordHash string
	OrgUserFk    sql.NullInt64
}
