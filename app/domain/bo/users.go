package bo

type User struct {
	Id           int
	Uuid         string
	PasswordHash string
	Email        string
	PhoneNumber  string
	UserName     string
	FirstName    string
	LastName     string
	ImageUrl     string
}
