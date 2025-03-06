package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/nilayrajderkar/redis-implementation/resp"
	"github.com/nilayrajderkar/redis-implementation/server"
)

func toString(resp interface{}) string {
	switch v := resp.(type) {
	case string:
		return v
	case *string:
		if v != nil {
			return *v
		}
	case []byte:
		return string(v)
	default:
		return fmt.Sprintf("%v", v)
	}
	return ""
}

func StartClient() error {
	// Connect to Redis server
	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		return fmt.Errorf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	fmt.Println("Connected to Redis server. Type your commands (press Ctrl+C to quit):")

	// Create a scanner to read user input
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		if input == "" {
			continue
		}

		// Custom parsing to handle quoted arguments
		var parts []string
		var currentPart strings.Builder
		inQuotes := false

		for i := 0; i < len(input); i++ {
			switch input[i] {
			case '"':
				inQuotes = !inQuotes
			case ' ':
				if !inQuotes {
					if currentPart.Len() > 0 {
						parts = append(parts, currentPart.String())
						currentPart.Reset()
					}
				} else {
					currentPart.WriteByte(input[i])
				}
			default:
				currentPart.WriteByte(input[i])
			}
		}

		if currentPart.Len() > 0 {
			parts = append(parts, currentPart.String())
		}

		if len(parts) == 0 {
			continue
		}

		// Convert each part into an array element, removing quotes if present
		args := make([]interface{}, len(parts))
		for i, part := range parts {
			// Remove surrounding quotes if present
			if len(part) >= 2 && part[0] == '"' && part[len(part)-1] == '"' {
				part = part[1 : len(part)-1]
			}
			args[i] = part
		}

		// Serialize the array of command and arguments
		serializedInput := resp.Serialize(args)

		deserializedResponse := server.HandleRequest(serializedInput)
		response, err := resp.Deserialize(deserializedResponse)
		if err != nil {
			fmt.Printf("Error deserializing response: %v\n", err)
			continue
		}

		fmt.Println(toString(response))
	}

	return nil
}
