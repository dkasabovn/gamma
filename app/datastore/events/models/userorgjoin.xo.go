package models

// Code generated by xo. DO NOT EDIT.

import (
	"context"
	"database/sql"
)

// UserOrgJoin represents a row from 'public.user_org_join'.
type UserOrgJoin struct {
	PermissionsCode  sql.NullInt64  `json:"PermissionsCode"`  // PermissionsCode
	Name             sql.NullString `json:"Name"`             // Name
	Description      sql.NullString `json:"Description"`      // Description
	OrganizationUUID sql.NullString `json:"OrganizationUuid"` // OrganizationUuid
}

// UserOrgJoinsByUserID runs a custom query, returning results as UserOrgJoin.
func UserOrgJoinsByUserID(ctx context.Context, db DB, userId int) ([]*UserOrgJoin, error) {
	// query
	const sqlstr = `SELECT ` +
		`"OrgUsers"."PermissionsCode"::INTEGER, ` +
		`"Organizations"."Name"::VARCHAR(20), ` +
		`"Organizations"."Description"::VARCHAR(255), ` +
		`"Organizations"."OrganizationUuid"::VARCHAR(36) ` +
		`FROM public."Users" ` +
		`INNER JOIN public."OrgUsers" ON "Users"."UserID" = "OrgUsers"."UserFk" ` +
		`INNER JOIN public."Organizations" ON "OrgUsers"."OrgFk" = "Organizations"."OrganizationID" ` +
		`WHERE "Users"."UserID" = $1`
	// run
	logf(sqlstr, userId)
	rows, err := db.QueryContext(ctx, sqlstr, userId)
	if err != nil {
		return nil, logerror(err)
	}
	defer rows.Close()
	// load results
	var res []*UserOrgJoin
	for rows.Next() {
		var uoj UserOrgJoin
		// scan
		if err := rows.Scan(&uoj.PermissionsCode, &uoj.Name, &uoj.Description, &uoj.OrganizationUUID); err != nil {
			return nil, logerror(err)
		}
		res = append(res, &uoj)
	}
	if err := rows.Err(); err != nil {
		return nil, logerror(err)
	}
	return res, nil
}
