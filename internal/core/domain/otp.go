package domain

// OTP (One-Time Password) struct with Otp and OtpHash fields.
type OTP struct {
	// Actual OTP value.
	Otp string `json:"otp"`

	// Hashed OTP value.
	OtpHash string `json:"otp_hash"`
}
