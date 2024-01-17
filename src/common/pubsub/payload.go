package pubsub

type Payload interface {
	GetTrackId() string
}
