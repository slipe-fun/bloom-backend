package service

func Strptr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
