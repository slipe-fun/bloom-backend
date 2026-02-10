package verification

type VerificationApp struct {
	verification VerificationRepo
}

func NewAuthApp(verification VerificationRepo) *VerificationApp {
	return &VerificationApp{
		verification: verification,
	}
}
