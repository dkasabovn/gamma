package auth

type UserSignup struct {
	Email       string `json:"email"`
	RawPassword string `json:"password"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	UserName    string `json:"user_name"`
}

type UserSignIn struct {
	Email       string `json:"email"`
	RawPassword string `json:"password"`
}
