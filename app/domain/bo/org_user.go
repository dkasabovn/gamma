package bo

type OrgUser struct {
	Id             int
	OrganizationFk int
	PoliciesNum    int
	UserFk         int
}

// TODO: Decide on perms
func (u OrgUser) CanCreateEvent() bool {
	return false
}
