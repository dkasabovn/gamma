package models

// Code generated by xo. DO NOT EDIT.

import (
	"context"
	"database/sql"
)

// Organization represents a row from 'public.Organizations'.
type Organization struct {
	OrganizationID   int            `json:"OrganizationID"`   // OrganizationID
	Name             sql.NullString `json:"Name"`             // Name
	Description      sql.NullString `json:"Description"`      // Description
	OrganizationUUID sql.NullString `json:"OrganizationUuid"` // OrganizationUuid
	// xo fields
	_exists, _deleted bool
}

// Exists returns true when the Organization exists in the database.
func (o *Organization) Exists() bool {
	return o._exists
}

// Deleted returns true when the Organization has been marked for deletion from
// the database.
func (o *Organization) Deleted() bool {
	return o._deleted
}

// Insert inserts the Organization to the database.
func (o *Organization) Insert(ctx context.Context, db DB) error {
	switch {
	case o._exists: // already exists
		return logerror(&ErrInsertFailed{ErrAlreadyExists})
	case o._deleted: // deleted
		return logerror(&ErrInsertFailed{ErrMarkedForDeletion})
	}
	// insert (primary key generated and returned by database)
	const sqlstr = `INSERT INTO public.Organizations (` +
		`Name, Description, OrganizationUuid` +
		`) VALUES (` +
		`$1, $2, $3` +
		`) RETURNING OrganizationID`
	// run
	logf(sqlstr, o.Name, o.Description, o.OrganizationUUID)
	if err := db.QueryRowContext(ctx, sqlstr, o.Name, o.Description, o.OrganizationUUID).Scan(&o.OrganizationID); err != nil {
		return logerror(err)
	}
	// set exists
	o._exists = true
	return nil
}

// Update updates a Organization in the database.
func (o *Organization) Update(ctx context.Context, db DB) error {
	switch {
	case !o._exists: // doesn't exist
		return logerror(&ErrUpdateFailed{ErrDoesNotExist})
	case o._deleted: // deleted
		return logerror(&ErrUpdateFailed{ErrMarkedForDeletion})
	}
	// update with composite primary key
	const sqlstr = `UPDATE public.Organizations SET ` +
		`Name = $1, Description = $2, OrganizationUuid = $3 ` +
		`WHERE OrganizationID = $4`
	// run
	logf(sqlstr, o.Name, o.Description, o.OrganizationUUID, o.OrganizationID)
	if _, err := db.ExecContext(ctx, sqlstr, o.Name, o.Description, o.OrganizationUUID, o.OrganizationID); err != nil {
		return logerror(err)
	}
	return nil
}

// Save saves the Organization to the database.
func (o *Organization) Save(ctx context.Context, db DB) error {
	if o.Exists() {
		return o.Update(ctx, db)
	}
	return o.Insert(ctx, db)
}

// Upsert performs an upsert for Organization.
func (o *Organization) Upsert(ctx context.Context, db DB) error {
	switch {
	case o._deleted: // deleted
		return logerror(&ErrUpsertFailed{ErrMarkedForDeletion})
	}
	// upsert
	const sqlstr = `INSERT INTO public.Organizations (` +
		`OrganizationID, Name, Description, OrganizationUuid` +
		`) VALUES (` +
		`$1, $2, $3, $4` +
		`)` +
		` ON CONFLICT (OrganizationID) DO ` +
		`UPDATE SET ` +
		`Name = EXCLUDED.Name, Description = EXCLUDED.Description, OrganizationUuid = EXCLUDED.OrganizationUuid `
	// run
	logf(sqlstr, o.OrganizationID, o.Name, o.Description, o.OrganizationUUID)
	if _, err := db.ExecContext(ctx, sqlstr, o.OrganizationID, o.Name, o.Description, o.OrganizationUUID); err != nil {
		return logerror(err)
	}
	// set exists
	o._exists = true
	return nil
}

// Delete deletes the Organization from the database.
func (o *Organization) Delete(ctx context.Context, db DB) error {
	switch {
	case !o._exists: // doesn't exist
		return nil
	case o._deleted: // deleted
		return nil
	}
	// delete with single primary key
	const sqlstr = `DELETE FROM public.Organizations ` +
		`WHERE OrganizationID = $1`
	// run
	logf(sqlstr, o.OrganizationID)
	if _, err := db.ExecContext(ctx, sqlstr, o.OrganizationID); err != nil {
		return logerror(err)
	}
	// set deleted
	o._deleted = true
	return nil
}

// OrganizationByOrganizationID retrieves a row from 'public.Organizations' as a Organization.
//
// Generated from index 'Organizations_pkey'.
func OrganizationByOrganizationID(ctx context.Context, db DB, organizationID int) (*Organization, error) {
	// query
	const sqlstr = `SELECT ` +
		`OrganizationID, Name, Description, OrganizationUuid ` +
		`FROM public.Organizations ` +
		`WHERE OrganizationID = $1`
	// run
	logf(sqlstr, organizationID)
	o := Organization{
		_exists: true,
	}
	if err := db.QueryRowContext(ctx, sqlstr, organizationID).Scan(&o.OrganizationID, &o.Name, &o.Description, &o.OrganizationUUID); err != nil {
		return nil, logerror(err)
	}
	return &o, nil
}
