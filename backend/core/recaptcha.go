package core

type ReCAPTCHAVerifier interface {
	Verify(response string, action string) (bool, error)
}
