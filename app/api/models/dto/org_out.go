package dto

import (
	userRepo "gamma/app/datastore/pg"

	"github.com/google/uuid"
)

type Organization struct {
	ID          uuid.UUID
	OrgName     string
	City        string
	OrgImageUrl string
}

type OrganizationWithPermissions struct {
	Organization
	PoliciesNum int
}

func ConvertOrganization(org *userRepo.Organization) *Organization {
	return &Organization{
		ID:          org.ID,
		OrgName:     org.OrgName,
		City:        org.City,
		OrgImageUrl: org.OrgImageUrl,
	}
}

func ConvertOrganizationWithPermissions(org *userRepo.GetUserOrganizationsRow) *OrganizationWithPermissions {
	return &OrganizationWithPermissions{
		Organization: Organization{
			ID:          org.OrganizationFk,
			OrgName:     org.OrgName,
			City:        org.City,
			OrgImageUrl: org.OrgImageUrl,
		},
		PoliciesNum: int(org.PoliciesNum),
	}
}

func ConvertOrganizationsWithPermissions(orgs []*userRepo.GetUserOrganizationsRow) []*OrganizationWithPermissions {
	ret := make([]*OrganizationWithPermissions, len(orgs))
	for i, org := range orgs {
		ret[i] = ConvertOrganizationWithPermissions(org)
	}
	return ret
}
