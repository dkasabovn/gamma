package bo

type PolicyNumber int

const (
	AFFILIATED          PolicyNumber = 1
	SCAN_QR                          = 1 << 1
	ACCEPT_APPLICATIONS              = 1 << 2
	MODIFY_EVENTS                    = 1 << 3
	CREATE_EVENTS                    = 1 << 4
	MODIFY_ORGANIZATION              = 1 << 5
	CREATE_ORGANIZATION              = 1 << 6
	MODIFY_ORG_USERS                 = 1 << 7
	SUPER_ADMIN                      = 1 << 8
	OWNER                            = 1<<9 | SUPER_ADMIN | MODIFY_ORG_USERS | MODIFY_EVENTS | MODIFY_ORGANIZATION | CREATE_EVENTS | CREATE_ORGANIZATION | ACCEPT_APPLICATIONS
)

func (p PolicyNumber) Is(num PolicyNumber) bool {
	return p&num == num
}

func (p PolicyNumber) Can(num PolicyNumber) bool {
	return p&num == num
}

func Create(nums ...PolicyNumber) int {
	num := 0
	for _, v := range nums {
		num |= int(v)
	}
	return num
}
