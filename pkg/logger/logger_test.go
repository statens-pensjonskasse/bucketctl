package logger

import (
	"github.com/pterm/pterm"
	"testing"
)

func TestPlainSprintf(t *testing.T) {
	tests := []struct {
		name   string
		pretty string
		args   []interface{}
		want   string
	}{
		{
			name:   "Empty",
			pretty: "",
			args:   []interface{}{},
			want:   "",
		},
		{
			name:   "Colours",
			pretty: "\x1b[31mHello, %s! You're number %d\x1b[0m",
			args:   []interface{}{"Go", 1},
			want:   "Hello, Go! You're number 1",
		},
		{
			name:   "Pterm",
			pretty: pterm.Blue("Blue ") + pterm.Bold.Sprintf("Bold"),
			args:   nil,
			want:   "Blue Bold",
		},
		{
			name:   "Remove emojis",
			pretty: "🔥👾é☄️🔠",
			args:   nil,
			want:   "é",
		},
		{
			name:   "ReplaceNonPrintableChars",
			pretty: "Hello,\t%s!\nYou're number %d",
			args:   []interface{}{"Go", 1},
			want:   "Hello,\tGo!\nYou're number 1",
		},
		{
			name:   "æøå🤷‍♀️ÆØÅ",
			pretty: "  æøåÆØÅ  ",
			args:   nil,
			want:   "æøåÆØÅ",
		},
		{
			name:   "TrimSpaces",
			pretty: " \t  Hello, %s! You're number %d    \t  ",
			args:   []interface{}{"Go", 9_999},
			want:   "Hello, Go! You're number 9999",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := plainSprintf(tt.pretty, tt.args...); got != tt.want {
				t.Errorf("got '%v', want '%v'", got, tt.want)
			}
		})
	}
}
