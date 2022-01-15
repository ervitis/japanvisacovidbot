package email

type (
	Properties struct {
		From    string
		To      string
		Headers map[string]string
	}

	IEmail interface {
		Send() error
		Properties() *Properties
	}
)
