package resp

import (
	"strconv"
)

func serializeString(s string) string {
	return "+" + s + "\r\n"
}

func serializeInteger(i int) string {
	return ":" + strconv.Itoa(i) + "\r\n"
}

func serializeBulkString(s string) string {
	return "$" + strconv.Itoa(len(s)) + "\r\n" + s + "\r\n"
}

func serializeError(s string) string {
	return "-" + s + "\r\n"
}

func serializeArray(elements []interface{}) string {
	bulkString := ""
	bulkString += "*" + strconv.Itoa(len(elements)) + "\r\n"
	for _, element := range elements {
		bulkString += Serialize(element)
	}
	return bulkString
}

func Serialize(element interface{}) string {
	switch element := element.(type) {
	case string:
		return serializeBulkString(element)
	case int:
		return serializeInteger(element)
	case []interface{}:
		return serializeArray(element)
	case error:
		return serializeError(element.Error())
	}
	return ""
}
