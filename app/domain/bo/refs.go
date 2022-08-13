package bo

func NewIntPtr(a int) *int {
	b := a
	return &b
}

func NewStrPtr(s string) *string {
	ps := s
	return &ps
}
