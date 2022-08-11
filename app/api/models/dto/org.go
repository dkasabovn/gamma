package dto

import userRepo "gamma/app/datastore/pg"

type ResOrganization struct {
	OrgName     string `json:"org_name"`
	City        string `json:"city"`
	Uuid        string `json:"uuid"`
	OrgImageUrl string `json:"org_image"`
	PolicyNum   int32  `json:"policy_num"`
}

func ConvertOrganization(org *userRepo.Organization) *ResOrganization {
	return &ResOrganization{
		OrgName:     org.OrgName,
		City:        org.City,
		Uuid:        org.Uuid,
		OrgImageUrl: org.OrgImageUrl,
	}
}

func ConvertUserOrganizationRow(org *userRepo.GetUserOrganizationsRow) *ResOrganization {
	return &ResOrganization{
		OrgName:     org.OrgName,
		City:        org.City,
		Uuid:        org.Uuid,
		OrgImageUrl: org.OrgImageUrl,
		PolicyNum:   org.PoliciesNum,
	}
}

func ConvertUserOrganizationRows(orgs []*userRepo.GetUserOrganizationsRow) []*ResOrganization {
	org_list := make([]*ResOrganization, len(orgs))
	for i, org := range orgs {
		org_list[i] = ConvertUserOrganizationRow(org)
	}
	return org_list
}
