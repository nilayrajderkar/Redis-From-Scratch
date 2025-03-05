package resp

import (
	"reflect"
	"testing"
)

func Test_deserializeSimpleString(t *testing.T) {
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
			name: "It should deserialize simple input string",
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
			got, err := deserializeString(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("deserialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("deserialize() returned nil, want %v", tt.want)
				return
			}
			if !tt.wantErr && *got != tt.want {
				t.Errorf("deserialize() = %v, want %v", *got, tt.want)
			}
		})
	}
}

func Test_deserializeInteger(t *testing.T) {
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
			got, err := deserializeInteger(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("deserialize_integer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("deserialize_integer() returned nil, want %v", tt.want)
				return
			}
			if !tt.wantErr && *got != tt.want {
				t.Errorf("deserialize_integer() = %v, want %v", *got, tt.want)
			}
		})
	}
}

func Test_deserializeBulkString(t *testing.T) {
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
			got, err := deserializeBulkString(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("deserialize_bulk_string() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("deserialize_bulk_string() returned nil, want %v", tt.want)
				return
			}
			if !tt.wantErr && *got != tt.want {
				t.Errorf("deserialize_bulk_string() = %v, want %v", *got, tt.want)
			}
		})
	}
}

func Test_deserializeError(t *testing.T) {
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
			got, err := deserializeError(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("deserialize_error() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("deserialize_error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_deserializeArray(t *testing.T) {
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
			got, err := deserializeArray(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("deserializeArray() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Errorf("deserializeArray() returned nil, want %v", tt.want)
					return
				}
				if len(*got) != len(*tt.want) {
					t.Errorf("deserializeArray() returned array of length %d, want %d", len(*got), len(*tt.want))
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
						t.Errorf("deserializeArray()[%d] = %v (%T), want %v (%T)",
							i, gotVal, gotVal, (*tt.want)[i], (*tt.want)[i])
					}
				}
			}
		})
	}
}

func Test_Deserialize(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    interface{}
		wantErr bool
	}{
		{
			name:    "should deserialize simple string",
			input:   "+OK\r\n",
			want:    "OK",
			wantErr: false,
		},
		{
			name:    "should deserialize integer",
			input:   ":123\r\n",
			want:    123,
			wantErr: false,
		},
		{
			name:    "should deserialize bulk string",
			input:   "$5\r\nhello\r\n",
			want:    "hello",
			wantErr: false,
		},
		{
			name:    "should deserialize error",
			input:   "-Error message\r\n",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "should deserialize array",
			input:   "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n",
			want:    []interface{}{"hello", "world"},
			wantErr: false,
		},
		{
			name:    "should deserialize mixed array",
			input:   "*3\r\n:1\r\n$5\r\nhello\r\n+OK\r\n",
			want:    []interface{}{1, "hello", "OK"},
			wantErr: false,
		},
		{
			name:    "should return error for invalid input",
			input:   "invalid",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Deserialize(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Deserialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Errorf("Deserialize() returned nil, want %v", tt.want)
					return
				}

				// Handle array type
				if gotArray, ok := got.(*[]interface{}); ok {
					derefArray := make([]interface{}, len(*gotArray))
					for i, v := range *gotArray {
						// Dereference each element if it's a pointer
						if intPtr, ok := v.(*int); ok {
							derefArray[i] = *intPtr
						} else if strPtr, ok := v.(*string); ok {
							derefArray[i] = *strPtr
						} else {
							derefArray[i] = v
						}
					}
					got = derefArray
				} else {
					// Handle single value pointer types
					if gotPtr, ok := got.(*int); ok {
						got = *gotPtr
					}
					if gotPtr, ok := got.(*string); ok {
						got = *gotPtr
					}
				}

				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Deserialize() = %v (%T), want %v (%T)", got, got, tt.want, tt.want)
				}
			}
		})
	}
}
