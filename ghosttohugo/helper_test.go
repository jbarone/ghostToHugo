package ghosttohugo

import (
	"encoding/json"
	"testing"
)

func Test_parseBool(t *testing.T) {
	tests := []struct {
		name string
		arg  json.RawMessage
		want bool
	}{
		{"nil", json.RawMessage(nil), false},
		{"empty", json.RawMessage([]byte{}), false},
		{"true", json.RawMessage([]byte("true")), true},
		{"false", json.RawMessage([]byte("false")), false},
		{"true_int", json.RawMessage([]byte{49}), true},
		{"false_int", json.RawMessage([]byte{0}), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseBool(tt.arg); got != tt.want {
				t.Errorf("parseBool() = %v, want %v", got, tt.want)
			}
		})
	}
}
