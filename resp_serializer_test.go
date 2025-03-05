package main

import (
	"reflect"
	"testing"
)

func Test_serializeSimpleString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "It should serialize simple input string",
			args: args{
				s: "+OK\r\n",
			},
			want:    "OK",
			wantErr: false,
		},
		{
			name: "It should raise error for incorrect simple input string",
			args: args{
				s: "+OK",
			},
			wantErr: true,
			want:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := serializeString(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("serialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("serialize() returned nil, want %v", tt.want)
				return
			}
			if !tt.wantErr && *got != tt.want {
				t.Errorf("serialize() = %v, want %v", *got, tt.want)
			}
		})
	}
}

func Test_serializeInteger(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "It should raise error for incorrect simple input integer",
			args: args{
				s: ":123",
			},
			wantErr: true,
		},
		{
			name: "It should return correct integer value",
			args: args{
				s: ":123\r\n",
			},
			want:    123,
			wantErr: false,
		},
		{

			name: "It should raise error if the input string is not a valid integer",
			args: args{
				s: ":123a\r\n",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := serializeInteger(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("serialize_integer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("serialize_integer() returned nil, want %v", tt.want)
				return
			}
			if !tt.wantErr && *got != tt.want {
				t.Errorf("serialize_integer() = %v, want %v", *got, tt.want)
			}
		})
	}
}

func Test_serializeBulkString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "It should raise error for incorrect bulk string input",
			args: args{
				s: "$123",
			},
			wantErr: true,
		},
		{
			name: "it should return correct bulk string value",
			args: args{
				s: "$5\r\nhello\r\n",
			},
			want:    "hello",
			wantErr: false,
		},
		{
			name: "it should raise error if the bulk string is too short",
			args: args{
				s: "$5\r\nhelo\r\n",
			},
			wantErr: true,
		},
		{
			name: "it should raise error if the bulk string is too long",
			args: args{
				s: "$5\r\nhellllo\r\n",
			},
			wantErr: true,
		},
		{
			name: "it should raise error if the bulk string is not properly terminated",
			args: args{
				s: "$5\r\nhello",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := serializeBulkString(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("serialize_bulk_string() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("serialize_bulk_string() returned nil, want %v", tt.want)
				return
			}
			if !tt.wantErr && *got != tt.want {
				t.Errorf("serialize_bulk_string() = %v, want %v", *got, tt.want)
			}
		})
	}
}

func Test_serializeError(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    *string
		wantErr bool
	}{
		{
			name: "it should return error message",
			args: args{
				s: "-ERR some error message\r\n",
			},
			wantErr: true,
		},
		{
			name: "it should return error message for incorrect input",
			args: args{
				s: "-ERR some error message",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := serializeError(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("serailize_error() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("serailize_error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_serializeArray(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    *[]interface{}
		wantErr bool
	}{
		{
			name: "it should return correct array value",
			args: args{
				s: "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n",
			},
			want:    &[]interface{}{"hello", "world"},
			wantErr: false,
		},
		{
			name: "it should return correct result for array of mixed types",
			args: args{
				s: "*3\r\n:1\r\n$5\r\nhello\r\n+OK\r\n",
			},
			want:    &[]interface{}{1, "hello", "OK"},
			wantErr: false,
		},
		{
			name: "it should handle nested arrays correctly",
			args: args{
				s: "*3\r\n:1\r\n*2\r\n+OK\r\n$5\r\nhello\r\n$5\r\nworld\r\n",
			},
			want:    &[]interface{}{1, []interface{}{"OK", "hello"}, "world"},
			wantErr: false,
		},
		{
			name: "it should raise error for incorrect array input",
			args: args{
				s: "*2\r\n$5\r\nhello\r\n$5\r\nworld",
			},
			wantErr: true,
		},
		{
			name: "it should raise error for incorrect array input",
			args: args{
				s: "*2\r\n$5\r\nhello\r\n$5\r\nworld",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := serializeArray(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("serializeArray() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Errorf("serializeArray() returned nil, want %v", tt.want)
					return
				}
				if len(*got) != len(*tt.want) {
					t.Errorf("serializeArray() returned array of length %d, want %d", len(*got), len(*tt.want))
					return
				}
				for i := range *got {
					gotVal := (*got)[i]
					if gotPtr, ok := gotVal.(*int); ok {
						gotVal = *gotPtr
					}
					if gotPtr, ok := gotVal.(*string); ok {
						gotVal = *gotPtr
					}

					if !reflect.DeepEqual(gotVal, (*tt.want)[i]) {
						t.Errorf("serializeArray()[%d] = %v (%T), want %v (%T)",
							i, gotVal, gotVal, (*tt.want)[i], (*tt.want)[i])
					}
				}
			}
		})
	}
}
