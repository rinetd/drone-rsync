package utils

import (
	"testing"
)

func TestReplace(t *testing.T) {
	type args struct {
		input     string
		escape    string
		delimiter string
		new       string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "字符测试",
			args: args{
				input:     "aa,bb",
				escape:    ",",
				delimiter: ",",
				new:       "\n",
			},
			want: "aa\nbb",
		},
		{name: "字符测试1",
			args: args{
				input:     "aa\t,bb",
				escape:    ",",
				delimiter: ",",
				new:       "\n",
			},
			want: "aa\t\nbb",
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Replace(tt.args.input, tt.args.escape, tt.args.delimiter, tt.args.new); got != tt.want {
				t.Errorf("Replace() = %v, want %v", got, tt.want)
			}
		})
	}
}
