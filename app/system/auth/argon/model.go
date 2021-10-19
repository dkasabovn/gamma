package argon

type UserAuth struct {
	UserIdentification string `json:"uuid"`
	Password           string `json:"password"`
}
