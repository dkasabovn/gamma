package dto

type UserSignUp struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	RawPassword string `json:"password"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	UserName    string `json:"username"`
}

type UserSignIn struct {
	Email       string `json:"email"`
	RawPassword string `json:"password"`
}

type UserUpdate struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	RawPassword string `json:"password"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	UserName    string `json:"username"`
	ImageUrl    string `json:"imageurl"`
}
