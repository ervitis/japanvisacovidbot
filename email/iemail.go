package email

type (
	IEmail interface {
		Send() error
	}
)
