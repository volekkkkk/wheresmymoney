package notification

type Provider interface {
	Send(message string) error
}
