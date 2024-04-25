package domain

type OTP struct {
	Otp     string `json:"otp"`
	OtpHash string `json:"otp_hash"`
}
