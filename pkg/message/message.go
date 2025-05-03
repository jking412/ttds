package message

type Manager interface {
	CreateChannel(name string) (chan string, error)
	RemoveChannel(name string) error
	SendMessage(channel string, message string) error
	GetChannel(channel string) (chan string, error)
}
