package bo

import "encoding/json"

type InviteTarget int

type InviteOption func(*Invite)

const (
	EVENT InviteTarget = iota
	ORGANIZATION
)

type Invite struct {
	Target       InviteTarget
	Requirements InviteRequirements
}

type InviteContext struct {
	OrganizationID string
	UserID         string
}

type InviteRequirements struct {
	OrganizationID *string
	UserID         *string
}

func NewInvite(opts ...InviteOption) *Invite {
	i := &Invite{
		Target:       EVENT,
		Requirements: InviteRequirements{},
	}

	for _, opt := range opts {
		opt(i)
	}

	return i
}

func NewInviteWithReqs(target InviteTarget, reqs InviteRequirements) *Invite {
	return &Invite{
		Target:       target,
		Requirements: reqs,
	}
}

func WithEvent() InviteOption {
	return func(i *Invite) {
		i.Target = EVENT
	}
}

func WithOrganization() InviteOption {
	return func(i *Invite) {
		i.Target = ORGANIZATION
	}
}

func WithOrgID(id string) InviteOption {
	return func(i *Invite) {
		i.Requirements.OrganizationID = NewStrPtr(id)
	}
}

func WithUserID(id string) InviteOption {
	return func(i *Invite) {
		i.Requirements.UserID = NewStrPtr(id)
	}
}

func (i *Invite) Validate(c *InviteContext) bool {
	if *i.Requirements.OrganizationID != c.OrganizationID && i.Requirements.OrganizationID != nil {
		return false
	}

	if *i.Requirements.UserID != c.UserID && i.Requirements.UserID != nil {
		return false
	}

	return true
}

func (i *Invite) Serialize() ([]byte, error) {
	return json.Marshal(i)
}
