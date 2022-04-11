package bo

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type InviteType string 
const (
	InviteToEvent = InviteType("event")
	InviteToOrg   = InviteType("organization")
)

type InviteConstraint string
const (
	ConstraintEmail = InviteConstraint("email")
	ConstraintOrg    = InviteConstraint("organization")
)

type InvitePolicy  struct {
	InviteType  	InviteType       `json:"invite_type"`
	InviteTo	    int              `json:"invite_to"`
	Constraint      InviteConstraint `json:"receiver_type"`
	Receiver        string	         `json:"receiver"`
	PoliciesNum     int              `json:"policies_num,omitempty"`
}

type Invite struct {
	Id				int
	Uuid 			string
	ExpirationDate  time.Time
	UseLimit		int
	Policy			InvitePolicy
}

func (p InvitePolicy) Value() (driver.Value, error) {
    return json.Marshal(p)
}

func (p *InvitePolicy) Scan(value interface{}) error {
	b, ok := value.([]byte)
    if !ok {
        return errors.New("type assertion to []byte failed")
    }

    return json.Unmarshal(b, &p)
}

func (i *Invite) UserCanAttend(user *User, userOrgs []OrganizationUserJoin) bool {

	if (i.UseLimit <= 0)  {
		return false
	}

	switch i.Policy.Constraint {
	case ConstraintEmail:
		return  user.Email == i.Policy.Receiver
	case ConstraintOrg:
		for _, orgJoin := range userOrgs {
			if orgJoin.Uuid == i.Policy.Receiver {
				return true
			}
		}
		return false

	default:
		return false
	}
}