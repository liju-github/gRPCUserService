package model

type User struct {
	ID             string
	Email          string
	PasswordHash   string
	Name           string
	StreetName     string
	Locality       string
	State          string
	Pincode        string
	PhoneNumber    string
	Reputation     int32
	VerificationCode string
	IsVerified     bool
}
