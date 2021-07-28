package ghosttohugo

import (
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
		Type string
		want bool
	}{
		{"nil", "", false},
		{"empty", "", false},
		{"true", "page", true},
		{"false", "false", false},
		{"true_int", "1", false},
		{"false_int", "0", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := post{
				Type: tt.Type,
			}
			if got := p.isPage(); got != tt.want {
				t.Errorf("post.isPage() = %v, want %v", got, tt.want)
			}
		})
	}
}
