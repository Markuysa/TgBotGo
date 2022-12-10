package events

type Fetcher interface {
	Fetch(limit int) ([]Event, error)
}
type Processor interface {
	Process(event Event) error
}

type Type int

const (
	Unknown Type = iota
	Message
)

type From struct {
	Username string `json:"username"`
}
type Chat struct {
	ID int `json:"id"`
}
type IncomingMessage struct {
	Text string `json:"text"`
	From From   `json:"from"`
	Chat Chat   `json:"chat"`
}

type Event struct {
	Type Type
	Text string
	Meta interface{}
}
