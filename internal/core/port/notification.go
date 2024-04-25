package port

type NotificationService interface {
	SendSMS(to string, template string) (bool, error)
	SendEmail(to string, template string) (bool, error)
}
