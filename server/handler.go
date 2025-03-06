package server

import (
	"errors"
	"strings"

	"github.com/nilayrajderkar/redis-implementation/resp"
)

// HandleRequest takes a serialized RESP string, deserializes it,
// processes the command, and returns a serialized response
func HandleRequest(input string) string {
	deserialized, err := resp.Deserialize(input)
	if err != nil {
		return resp.Serialize(err)
	}

	array, ok := deserialized.(*[]interface{})
	if !ok {
		return resp.Serialize(errors.New("invalid command format"))
	}

	// Need at least one element (the command)
	if len(*array) == 0 {
		return resp.Serialize(errors.New("empty command"))
	}

	// First element should be the command string
	command, ok := (*array)[0].(*string)
	if !ok {
		return resp.Serialize(errors.New("command must be a string"))
	}

	// Convert command to uppercase for case-insensitive comparison
	switch strings.ToUpper(*command) {
	case "PING":
		return resp.Serialize("PONG")
	case "ECHO":
		if len(*array) < 2 {
			return resp.Serialize(errors.New("ECHO requires an argument"))
		}
		// Echo back the second element
		if str, ok := (*array)[1].(*string); ok {
			return resp.Serialize(*str)
		}
		return resp.Serialize(errors.New("ECHO argument must be a string"))
	default:
		return resp.Serialize(errors.New("Unknown command '" + *command + "'"))
	}
}
