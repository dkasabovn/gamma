package bo

type Organization struct {
	Uuid             string
	Id               int
	OrganizationName string
	City             string
}

type OrganizationUser struct {
	Uuid             string
	Id               int
	OrganizationName string
	City             string
	PolicyNum        int
}
