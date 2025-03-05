package resp

import (
	"errors"
	"strconv"
	"strings"
)

func deserializeString(s string) (*string, error) {
	index := strings.Index(s, "\r\n")
	if index == -1 {
		return nil, errors.New("invalid input string")
	}
	result := s[1:index]
	return &result, nil
}

func deserializeInteger(s string) (*int, error) {
	index := strings.Index(s, "\r\n")
	if index == -1 {
		return nil, errors.New("invalid input string for integer type")
	}

	integerValue, err := strconv.Atoi(s[1:index])
	if err != nil {
		return nil, errors.New("cannot convert to integer, invalid value")
	}
	return &integerValue, nil
}

func deserializeBulkString(s string) (*string, error) {
	index := strings.Index(s, "\r\n")
	if index == -1 {
		return nil, errors.New("invalid input string for bulk string type")
	}
	// first character is the length of the string
	length, err := strconv.Atoi(s[1:index])
	if err != nil {
		return nil, errors.New("invalid length for bulk string")
	}
	if length == -1 {
		return nil, errors.New("key does not exist")
	}

	// Check for second \r\n for the bulk string content
	expectedEnd := index + 2 + length + 2
	if expectedEnd > len(s) {
		return nil, errors.New("bulk string is too short")
	}
	if s[index+2+length:expectedEnd] != "\r\n" {
		return nil, errors.New("bulk string is not properly terminated")
	}
	result := s[index+2 : index+2+length]
	return &result, nil
}

func deserializeError(s string) (*string, error) {
	index := strings.Index(s, "\r\n")
	if index == -1 {
		return nil, errors.New("invalid input string for error type")
	}
	return nil, errors.New(s[1:index])
}

func deserializeArray(s string) (*[]interface{}, error) {
	index := strings.Index(s, "\r\n")
	if index == -1 {
		return nil, errors.New("invalid input string for array type")
	}
	numberOfElements, err := strconv.Atoi(s[1:index])
	if err != nil {
		return nil, errors.New("invalid number of elements for array")
	}
	result := make([]interface{}, 0, numberOfElements)
	currentPos := index + 2

	for i := 0; i < numberOfElements; i++ {
		if currentPos >= len(s) {
			return nil, errors.New("unexpected end of array")
		}

		// Process each element based on its type
		elementType := s[currentPos]
		var elementEnd int

		switch elementType {
		case '$': // Bulk String
			// Find length
			firstNewline := strings.Index(s[currentPos:], "\r\n")
			if firstNewline == -1 {
				return nil, errors.New("invalid bulk string format in array")
			}
			length, err := strconv.Atoi(s[currentPos+1 : currentPos+firstNewline])
			if err != nil {
				return nil, errors.New("invalid bulk string length in array")
			}
			elementEnd = currentPos + firstNewline + 2 + length + 2 // Include both \r\n
		case ':': // Integer
			nextNewline := strings.Index(s[currentPos:], "\r\n")
			if nextNewline == -1 {
				return nil, errors.New("invalid integer format in array")
			}
			elementEnd = currentPos + nextNewline + 2
		case '+': // Simple String
			nextNewline := strings.Index(s[currentPos:], "\r\n")
			if nextNewline == -1 {
				return nil, errors.New("invalid simple string format in array")
			}
			elementEnd = currentPos + nextNewline + 2
		case '-': // Error
			nextNewline := strings.Index(s[currentPos:], "\r\n")
			if nextNewline == -1 {
				return nil, errors.New("invalid error format in array")
			}
			elementEnd = currentPos + nextNewline + 2
		case '*': // Nested Array
			// Find the complete nested array by counting elements
			nestedNewline := strings.Index(s[currentPos:], "\r\n")
			if nestedNewline == -1 {
				return nil, errors.New("invalid nested array format")
			}
			nestedCount, err := strconv.Atoi(s[currentPos+1 : currentPos+nestedNewline])
			if err != nil {
				return nil, errors.New("invalid nested array length")
			}
			// Recursively process the nested array
			nested, err := Deserialize(s[currentPos:])
			if err != nil {
				return nil, err
			}
			result = append(result, nested)
			// Update elementEnd based on where the nested array processing ended
			elementEnd = len(s)
			for tmp := currentPos; tmp < len(s); tmp++ {
				if nestedCount == 0 {
					elementEnd = tmp
					break
				}
				if strings.HasPrefix(s[tmp:], "\r\n") {
					nestedCount--
				}
			}
		default:
			return nil, errors.New("invalid element type in array")
		}

		if elementEnd > len(s) {
			return nil, errors.New("unexpected end of array element")
		}

		element, err := Deserialize(s[currentPos:elementEnd])
		if err != nil {
			return nil, err
		}

		result = append(result, element)
		currentPos = elementEnd
	}

	return &result, nil
}

func Deserialize(s string) (interface{}, error) {
	if len(s) == 0 {
		return nil, errors.New("empty input string")
	}

	switch s[0] {
	case '*':
		return deserializeArray(s)
	case '$':
		return deserializeBulkString(s)
	case ':':
		return deserializeInteger(s)
	case '-':
		return deserializeError(s)
	case '+':
		return deserializeString(s)
	default:
		return nil, errors.New("invalid input string")
	}
}
