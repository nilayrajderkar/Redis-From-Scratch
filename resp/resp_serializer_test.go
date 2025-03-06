package resp

import (
	"errors"
	"testing"
)

func Test_serializeString(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{
			name: "It should serialize simple string",
			arg:  "OK",
			want: "+OK\r\n",
		},
		{
			name: "It should serialize empty string",
			arg:  "",
			want: "+\r\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := serializeString(tt.arg)
			if got != tt.want {
				t.Errorf("serializeString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_serializeInteger(t *testing.T) {
	tests := []struct {
		name string
		arg  int
		want string
	}{
		{
			name: "It should serialize positive integer",
			arg:  123,
			want: ":123\r\n",
		},
		{
			name: "It should serialize zero",
			arg:  0,
			want: ":0\r\n",
		},
		{
			name: "It should serialize negative integer",
			arg:  -123,
			want: ":-123\r\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := serializeInteger(tt.arg)
			if got != tt.want {
				t.Errorf("serializeInteger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_serializeBulkString(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{
			name: "It should serialize bulk string",
			arg:  "hello",
			want: "$5\r\nhello",
		},
		{
			name: "It should serialize empty bulk string",
			arg:  "",
			want: "$0\r\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := serializeBulkString(tt.arg)
			if got != tt.want {
				t.Errorf("serializeBulkString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_serializeError(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{
			name: "It should serialize error message",
			arg:  "Error occurred",
			want: "-Error occurred\r\n",
		},
		{
			name: "It should serialize empty error message",
			arg:  "",
			want: "-\r\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := serializeError(tt.arg)
			if got != tt.want {
				t.Errorf("serializeError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_serializeArray(t *testing.T) {
	tests := []struct {
		name string
		arg  []interface{}
		want string
	}{
		{
			name: "It should serialize array of strings",
			arg:  []interface{}{"hello", "world"},
			want: "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n",
		},
		{
			name: "It should serialize array of mixed types",
			arg:  []interface{}{1, "hello", errors.New("error")},
			want: "*3\r\n:1\r\n$5\r\nhello\r\n-error\r\n",
		},
		{
			name: "It should serialize empty array",
			arg:  []interface{}{},
			want: "*0\r\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := serializeArray(tt.arg)
			if got != tt.want {
				t.Errorf("serializeArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Serialize(t *testing.T) {
	tests := []struct {
		name string
		arg  interface{}
		want string
	}{
		{
			name: "It should serialize string",
			arg:  "hello",
			want: "$5\r\nhello\r\n",
		},
		{
			name: "It should serialize integer",
			arg:  123,
			want: ":123\r\n",
		},
		{
			name: "It should serialize array",
			arg:  []interface{}{"hello", 123},
			want: "*2\r\n$5\r\nhello\r\n:123\r\n",
		},
		{
			name: "It should serialize error",
			arg:  errors.New("error message"),
			want: "-error message\r\n",
		},
		{
			name: "It should return empty string for unsupported type",
			arg:  3.14,
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Serialize(tt.arg)
			if got != tt.want {
				t.Errorf("Serialize() = %v, want %v", got, tt.want)
			}
		})
	}
}
