package command

type ChatCmd int

const (
	Ping              ChatCmd = 0
	Pong              ChatCmd = 10000
	Connect           ChatCmd = 100
	SendChat          ChatCmd = 3101
	RequestRecentChat ChatCmd = 5101
	Chat              ChatCmd = 93101
	Donation          ChatCmd = 93102
)

type PongResponse struct {
	Cmd ChatCmd `json:"cmd"`
	Ver string  `json:"ver"`
}

var PongInstance = PongResponse{Cmd: Pong, Ver: "3"}
