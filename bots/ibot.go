package bots

type (
	IBot interface {
		SendNotification(interface{}) error
		StartServer() error
		Close()
	}
)
