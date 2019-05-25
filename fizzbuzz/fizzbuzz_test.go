package fizzbuzz

import "testing"

func TestFizzBuzz(t *testing.T) {
	fizzBuzzTests := []struct {
		number int
		result string
	}{
		{1, "1"},
		{2, "2"},
		{3, "fizz"},
		{5, "buzz"},
		{15, "fizzbuzz"},
		{13, "fizz"},
		{58, "buzz"},
		{35, "fizzbuzz"},
	}

	for _, tt := range fizzBuzzTests {
		got := FizzBuzz(tt.number)
		want := tt.result
		if got != want {
			t.Errorf("%d got %s, want %s", tt.number, got, want)
		}
	}
}
