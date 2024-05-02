package port

// INotificationService interface defines methods for sending notifications (SMS, Email, Push, Whatsapp).
type INotificationService interface {
	// SendSMS sends an SMS to the specified recipient using the given template.
	SendSMS(to string, template string) (bool, error)

	// SendEmail sends an email to the specified recipient using the given template.
	SendEmail(to string, template string) (bool, error)

	// SendPush sends a push notification to the specified recipient using the given template.
	SendPush(to string, template string) (bool, error)

	// SendWhatsapp sends a Whatsapp message to the specified recipient using the given template.
	SendWhatsapp(to string, template string) (bool, error)
}
