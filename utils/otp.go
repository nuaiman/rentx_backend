package utils

import (
	"errors"
	"time"
)

// VerifyOTP verifies the OTP for a given phone number.
// In a real app, you'd check a DB/cache or external service.
func VerifyOTP(phone string, otp string) error {
	// TODO: Implement your OTP verification logic here:
	// - fetch the stored OTP for this phone from DB/cache
	// - check if it matches and is not expired
	// - return nil if valid
	// - return errors.New("invalid or expired OTP") if not valid

	// Example dummy logic for now:
	if otp == "123456" { // test OTP
		return nil
	}
	return errors.New("invalid or expired OTP")
}

// SendOTP generates and sends an OTP to the given phone.
func SendOTP(phone string) (string, error) {
	// TODO: generate OTP, store in DB/cache with expiration,
	// and send it via SMS provider
	otp := "123456" // hardcoded for testing
	_ = time.Now()  // simulate storing with expiration
	return otp, nil
}
