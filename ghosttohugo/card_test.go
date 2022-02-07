package ghosttohugo

import (
	"testing"
)

func Test_cardCode(t *testing.T) {
	type args struct {
		payload interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"non_map", args{nil}, ""},
		{"empty", args{map[string]interface{}{}}, ""},
		{
			"code",
			args{map[string]interface{}{
				"code": "def hello():\n    print('hello')",
			}},
			"```\ndef hello():\n    print('hello')\n```\n",
		},
		{
			"code_language",
			args{map[string]interface{}{
				"code":     "def hello():\n    print('hello')",
				"language": "python",
			}},
			"```python\ndef hello():\n    print('hello')\n```\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cardCode(tt.args.payload); got != tt.want {
				t.Errorf("cardCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cardEmbed(t *testing.T) {
	type args struct {
		payload interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"non_map", args{nil}, ""},
		{"empty", args{map[string]interface{}{}}, ""},
		{"empty", args{map[string]interface{}{"html": "test"}}, "test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cardEmbed(tt.args.payload); got != tt.want {
				t.Errorf("cardEmbed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cardHR(t *testing.T) {
	type args struct {
		payload interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"valid", args{nil}, "---\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cardHR(tt.args.payload); got != tt.want {
				t.Errorf("cardHR() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cardHTML(t *testing.T) {
	type args struct {
		payload interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"non_map", args{nil}, ""},
		{"empty", args{map[string]interface{}{}}, ""},
		{
			"valid",
			args{map[string]interface{}{"html": "<h1>test</h1>"}},
			"<h1>test</h1>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cardHTML(tt.args.payload); got != tt.want {
				t.Errorf("cardHTML() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cardImage(t *testing.T) {
	type args struct {
		payload interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"non_map", args{nil}, ""},
		{"empty", args{map[string]interface{}{}}, ""},
		{
			"src_no_caption",
			args{map[string]interface{}{
				"src": "test",
			}},
			"{{< figure src=\"test\" >}}\n",
		},
		{
			"src_caption",
			args{map[string]interface{}{
				"src":     "test",
				"caption": "caption",
			}},
			"{{< figure src=\"test\" caption=\"caption\" >}}\n",
		},
	}

	NewImgDownloader("", "", false)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cardImage(tt.args.payload); got != tt.want {
				t.Errorf("cardImage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cardMarkdown(t *testing.T) {
	type args struct {
		payload interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"non_map", args{nil}, ""},
		{"empty", args{map[string]interface{}{}}, ""},
		{
			"valid",
			args{map[string]interface{}{"markdown": "test"}},
			"test\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cardMarkdown(tt.args.payload); got != tt.want {
				t.Errorf("cardMarkdown() = %v, want %v", got, tt.want)
			}
		})
	}
}
