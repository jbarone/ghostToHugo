package ghosttohugo

import "testing"

func Test_atomSoftReturn(t *testing.T) {
	type args struct {
		value   string
		payload interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"softReturn", args{"", nil}, "\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := atomSoftReturn(tt.args.value, tt.args.payload); got != tt.want {
				t.Errorf("atomSoftReturn() = %v, want %v", got, tt.want)
			}
		})
	}
}
