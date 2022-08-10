package mqtt

import "fmt"

// nicked from github.com/goiiot/libmqtt
const (
	// MQTT 3.1.1 ConnAck code
	CodeUnacceptableVersion   = 1 // Packet: ConnAck
	CodeIdentifierRejected    = 2 // Packet: ConnAck
	CodeServerUnavailable     = 3 // Packet: ConnAck
	CodeBadUsernameOrPassword = 4 // Packet: ConnAck
	CodeUnauthorized          = 5 // Packet: ConnAck
)

func ReasonString(code byte) string {
	switch code {
	case CodeUnacceptableVersion:
		return "unacceptable version"
	case CodeIdentifierRejected:
		return "identifier rejected"
	case CodeServerUnavailable:
		return "server unavailable"
	case CodeBadUsernameOrPassword:
		return "bad user or password"
	case CodeUnauthorized:
		return "unauthorized"
	default:
		return fmt.Sprintf("unknown code %d", code)
	}
}
