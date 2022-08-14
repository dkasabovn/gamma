package bo

import "github.com/google/uuid"

type PartialUser struct {
	UUID     uuid.UUID
	Username string
}
