package notification

type SendEmailInput struct {
	Email   string
	Subject string
	Body    string
}

type ProcessingCompletedInput struct {
	VideoID  string
	Email    string
	Filename string
	Status   string
}
