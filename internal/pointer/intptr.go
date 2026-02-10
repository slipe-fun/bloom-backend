package pointer

func Intptr(i int) *int {
	if i == 0 {
		return nil
	}
	return &i
}
