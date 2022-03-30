package dto

type UserSignup struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	RawPassword string `json:"password"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	UserName    string `json:"user_name"`
	ImageUrl    string `json:"image_url"`
}

type UserSignIn struct {
	Email       string `json:"email"`
	RawPassword string `json:"password"`
}
