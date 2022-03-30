package bo

type Organization struct {
	Id               int
	Uuid             string
	OrganizationName string
	City             string
	ImageUrl         string
}

type OrganizationUserJoin struct {
	Organization
	PolicyNum int
}
