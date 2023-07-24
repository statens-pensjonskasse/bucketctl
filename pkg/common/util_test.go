package common

import "testing"

func TestSlicesContainsSameElements(t *testing.T) {
	type args[T comparable] struct {
		a []T
		b []T
	}
	type testCase[T comparable] struct {
		name string
		args args[T]
		want bool
	}
	stringTests := []testCase[string]{
		{
			name: "Different length",
			args: args[string]{
				a: []string{"x", "x"},
				b: []string{"x"},
			},
			want: false,
		},
		{
			name: "Different content",
			args: args[string]{
				a: []string{"x"},
				b: []string{"y"},
			},
			want: false,
		},
		{
			name: "Same elements, different amount of occurrences",
			args: args[string]{
				a: []string{"x", "x"},
				b: []string{"x"},
			},
			want: false,
		},
		{
			name: "Same content, same ordering",
			args: args[string]{
				a: []string{"x", "y"},
				b: []string{"x", "y"},
			},
			want: true,
		},
		{
			name: "Same content, different ordering",
			args: args[string]{
				a: []string{"x", "y", "z", "x"},
				b: []string{"y", "z", "x", "x"},
			},
			want: true,
		},
	}
	intTests := []testCase[int]{
		{
			name: "Same content, different ordering",
			args: args[int]{
				a: []int{0, 1, 2, 3, 4, 3, 3, 3},
				b: []int{3, 2, 4, 1, 0, 3, 3, 3},
			},
			want: true,
		},
	}
	for _, tt := range stringTests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SlicesContainsSameElements(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("SlicesContainsSameElements() = %v, want %v", got, tt.want)
			}
		})
	}
	for _, tt := range intTests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SlicesContainsSameElements(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("SlicesContainsSameElements() = %v, want %v", got, tt.want)
			}
		})
	}
}
