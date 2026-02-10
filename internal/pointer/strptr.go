package pointer

func Strptr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
