package ghosttohugo

import (
	"encoding/json"
	"testing"
)

func Test_post_isDraft(t *testing.T) {
	tests := []struct {
		name   string
		status string
		want   bool
	}{
		{"yes", "draft", true},
		{"no", "published", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := post{
				Status: tt.status,
			}
			if got := p.isDraft(); got != tt.want {
				t.Errorf("post.isDraft() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_post_isPage(t *testing.T) {
	tests := []struct {
		name string
		page json.RawMessage
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
			p := post{
				Page: tt.page,
			}
			if got := p.isPage(); got != tt.want {
				t.Errorf("post.isPage() = %v, want %v", got, tt.want)
			}
		})
	}
}
