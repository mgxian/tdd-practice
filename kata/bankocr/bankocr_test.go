package bankocr

import (
	"reflect"
	"testing"
)

func TestParseEntry(t *testing.T) {
	tests := []struct {
		entry string
		want  int
	}{
		{
			`
   
  |
  |
`, 1},
		{
			`
 _ 
 _|
|_ 
`, 2},
		{
			`
 _ 
 _|
 _|
`, 3},
	}

	for _, tt := range tests {
		want := tt.want
		got := parseEntry(tt.entry)
		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	}
}

func TestParseStringLine(t *testing.T) {
	aStringLine := `
 _     _  _     _  _  _  _  _ 
| |  | _| _||_||_ |_   ||_||_|
|_|  ||_  _|  | _||_|  ||_| _|	
`
	got := parseStringLine(aStringLine)
	want := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestParseNumbersFromFile(t *testing.T) {
	aFilePath := "./0123456789.txt"
	got := parseNumbersFromFile(aFilePath)
	want := [][]int{
		{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		{9, 8, 7, 6, 5, 4, 3, 2, 1, 0},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestValideAccountNumbers(t *testing.T) {
	testCases := []struct {
		accountNumbers []int
		isValid        bool
	}{
		{[]int{3, 4, 5, 8, 8, 2, 8, 6, 5}, true},
		{[]int{3, 1, 5, 8, 8, 2, 8, 6, 5}, false},
	}
	for _, tt := range testCases {
		got := validAccountNumbers(tt.accountNumbers)
		want := tt.isValid
		if got != want {
			t.Errorf("got %t, want %t", got, want)
		}
	}
}
