package models

// Code generated by xo. DO NOT EDIT.

import (
	"context"
	"database/sql"
)

// OrgUser represents a row from 'public.OrgUsers'.
type OrgUser struct {
	OrgUserID       int           `json:"OrgUserID"`       // OrgUserID
	PermissionsCode sql.NullInt64 `json:"PermissionsCode"` // PermissionsCode
	UserFk          sql.NullInt64 `json:"UserFk"`          // UserFk
	OrgFk           sql.NullInt64 `json:"OrgFk"`           // OrgFk
	// xo fields
	_exists, _deleted bool
}

// Exists returns true when the OrgUser exists in the database.
func (ou *OrgUser) Exists() bool {
	return ou._exists
}

// Deleted returns true when the OrgUser has been marked for deletion from
// the database.
func (ou *OrgUser) Deleted() bool {
	return ou._deleted
}

// Insert inserts the OrgUser to the database.
func (ou *OrgUser) Insert(ctx context.Context, db DB) error {
	switch {
	case ou._exists: // already exists
		return logerror(&ErrInsertFailed{ErrAlreadyExists})
	case ou._deleted: // deleted
		return logerror(&ErrInsertFailed{ErrMarkedForDeletion})
	}
	// insert (primary key generated and returned by database)
	const sqlstr = `INSERT INTO public.OrgUsers (` +
		`PermissionsCode, UserFk, OrgFk` +
		`) VALUES (` +
		`$1, $2, $3` +
		`) RETURNING OrgUserID`
	// run
	logf(sqlstr, ou.PermissionsCode, ou.UserFk, ou.OrgFk)
	if err := db.QueryRowContext(ctx, sqlstr, ou.PermissionsCode, ou.UserFk, ou.OrgFk).Scan(&ou.OrgUserID); err != nil {
		return logerror(err)
	}
	// set exists
	ou._exists = true
	return nil
}

// Update updates a OrgUser in the database.
func (ou *OrgUser) Update(ctx context.Context, db DB) error {
	switch {
	case !ou._exists: // doesn't exist
		return logerror(&ErrUpdateFailed{ErrDoesNotExist})
	case ou._deleted: // deleted
		return logerror(&ErrUpdateFailed{ErrMarkedForDeletion})
	}
	// update with composite primary key
	const sqlstr = `UPDATE public.OrgUsers SET ` +
		`PermissionsCode = $1, UserFk = $2, OrgFk = $3 ` +
		`WHERE OrgUserID = $4`
	// run
	logf(sqlstr, ou.PermissionsCode, ou.UserFk, ou.OrgFk, ou.OrgUserID)
	if _, err := db.ExecContext(ctx, sqlstr, ou.PermissionsCode, ou.UserFk, ou.OrgFk, ou.OrgUserID); err != nil {
		return logerror(err)
	}
	return nil
}

// Save saves the OrgUser to the database.
func (ou *OrgUser) Save(ctx context.Context, db DB) error {
	if ou.Exists() {
		return ou.Update(ctx, db)
	}
	return ou.Insert(ctx, db)
}

// Upsert performs an upsert for OrgUser.
func (ou *OrgUser) Upsert(ctx context.Context, db DB) error {
	switch {
	case ou._deleted: // deleted
		return logerror(&ErrUpsertFailed{ErrMarkedForDeletion})
	}
	// upsert
	const sqlstr = `INSERT INTO public.OrgUsers (` +
		`OrgUserID, PermissionsCode, UserFk, OrgFk` +
		`) VALUES (` +
		`$1, $2, $3, $4` +
		`)` +
		` ON CONFLICT (OrgUserID) DO ` +
		`UPDATE SET ` +
		`PermissionsCode = EXCLUDED.PermissionsCode, UserFk = EXCLUDED.UserFk, OrgFk = EXCLUDED.OrgFk `
	// run
	logf(sqlstr, ou.OrgUserID, ou.PermissionsCode, ou.UserFk, ou.OrgFk)
	if _, err := db.ExecContext(ctx, sqlstr, ou.OrgUserID, ou.PermissionsCode, ou.UserFk, ou.OrgFk); err != nil {
		return logerror(err)
	}
	// set exists
	ou._exists = true
	return nil
}

// Delete deletes the OrgUser from the database.
func (ou *OrgUser) Delete(ctx context.Context, db DB) error {
	switch {
	case !ou._exists: // doesn't exist
		return nil
	case ou._deleted: // deleted
		return nil
	}
	// delete with single primary key
	const sqlstr = `DELETE FROM public.OrgUsers ` +
		`WHERE OrgUserID = $1`
	// run
	logf(sqlstr, ou.OrgUserID)
	if _, err := db.ExecContext(ctx, sqlstr, ou.OrgUserID); err != nil {
		return logerror(err)
	}
	// set deleted
	ou._deleted = true
	return nil
}

// OrgUserByOrgUserID retrieves a row from 'public.OrgUsers' as a OrgUser.
//
// Generated from index 'OrgUsers_pkey'.
func OrgUserByOrgUserID(ctx context.Context, db DB, orgUserID int) (*OrgUser, error) {
	// query
	const sqlstr = `SELECT ` +
		`OrgUserID, PermissionsCode, UserFk, OrgFk ` +
		`FROM public.OrgUsers ` +
		`WHERE OrgUserID = $1`
	// run
	logf(sqlstr, orgUserID)
	ou := OrgUser{
		_exists: true,
	}
	if err := db.QueryRowContext(ctx, sqlstr, orgUserID).Scan(&ou.OrgUserID, &ou.PermissionsCode, &ou.UserFk, &ou.OrgFk); err != nil {
		return nil, logerror(err)
	}
	return &ou, nil
}

// Organization returns the Organization associated with the OrgUser's (OrgFk).
//
// Generated from foreign key 'OrgUsers_OrgFk_fkey'.
func (ou *OrgUser) Organization(ctx context.Context, db DB) (*Organization, error) {
	return OrganizationByOrganizationID(ctx, db, int(ou.OrgFk.Int64))
}

// User returns the User associated with the OrgUser's (UserFk).
//
// Generated from foreign key 'OrgUsers_UserFk_fkey'.
func (ou *OrgUser) User(ctx context.Context, db DB) (*User, error) {
	return UserByUserID(ctx, db, int(ou.UserFk.Int64))
}
