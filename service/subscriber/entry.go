package subscriber

type Entry struct {
	UserID string
	Email  string

	Channel int
}

const (
	ChannelInternal = iota
)
