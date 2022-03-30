package dto

import "gamma/app/domain/bo"

type EventByOrg struct {
	bo.Event `json:"event"`
	OrgUuid  string `json:"organization_uuid"`
}
