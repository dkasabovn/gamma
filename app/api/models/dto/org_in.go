package dto

type OrganizationGet struct {
	OrganizationID string `param:"organization_id"`
}

type OrganizationCreate struct {
	OrgName     string `json:"org_name"`
	City        string `json:"city"`
	OrgImageUrl string `json:"org_image_url"`
	AdminId     string `json:"admin_id"`
}
