package bo

type OrgUser struct {
	Id             int
	OrganizationFk int
	PoliciesNum    int
	UserFk         int
}
