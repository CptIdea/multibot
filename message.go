package multibot

type Message struct {
	Text    string
	FromID  int64
	PeerID  int64
	Payload string
}
