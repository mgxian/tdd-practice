package bankocr

import (
	"reflect"
	"sort"
	"strings"
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
	aFilePath := "./use-case-1.txt"
	got := parseNumbersFromFile(aFilePath)
	want := [][]int{
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{1, 1, 1, 1, 1, 1, 1, 1, 1},
		{2, 2, 2, 2, 2, 2, 2, 2, 2},
		{3, 3, 3, 3, 3, 3, 3, 3, 3},
		{4, 4, 4, 4, 4, 4, 4, 4, 4},
		{5, 5, 5, 5, 5, 5, 5, 5, 5},
		{6, 6, 6, 6, 6, 6, 6, 6, 6},
		{7, 7, 7, 7, 7, 7, 7, 7, 7},
		{8, 8, 8, 8, 8, 8, 8, 8, 8},
		{9, 9, 9, 9, 9, 9, 9, 9, 9},
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

func TestParseAndOutputEntry(t *testing.T) {
	output := new(strings.Builder)
	aFilePath := "./use-case-3.txt"
	parseAndOutputEntry(aFilePath, output)
	got := output.String()
	want := "000000051\n49006771? ILL\n1234?678? ILL\n"
	if got != want {
		t.Errorf("got '%s', want '%s'", got, want)
	}
}

func TestSmartParseEntry(t *testing.T) {
	tests := []struct {
		entry string
		want  []int
	}{
		{
			`
   
  |
  |
`, []int{1, 7}},
		{
			`
 _ 
 _|
|_ 
`, []int{2}},
		{
			`
 _ 
 _|
 _|
`, []int{3, 9}},
		{
			`
 _ 
|_|
|_|
`, []int{0, 8, 9, 6}},
	}

	for _, tt := range tests {
		want := tt.want
		got := smartParseEntry(tt.entry)
		assertSliceEqual(t, got, want)
	}
}

func TestSmartParseAndOutputEntry(t *testing.T) {
	output := new(strings.Builder)
	aFilePath := "./use-case-4.txt"
	smartParseAndOutputEntry(aFilePath, output)
	got := output.String()
	want := `711111111
777777177
200800000
333393333
888888888 AMB ['888886888', '888888988', '888888880']
555555555 AMB ['559555555', '555655555']
666666666 AMB ['686666666', '666566666']
999999999 AMB ['899999999', '993999999', '999959999']
490067715 AMB ['490867715', '490067719', '490067115']
123456789
000000051
490867715
`
	gotEntry := strings.Split(got, "\n")
	wantEntry := strings.Split(want, "\n")
	if len(gotEntry) != len(wantEntry) {
		t.Fatalf("got %d, want %d", len(gotEntry), len(wantEntry))
	}

	for i := 0; i < len(gotEntry); i++ {
		assertEntryResult(t, gotEntry[i], wantEntry[i])
	}
}

func assertEntryResult(t *testing.T, gotEntry, wantEntry string) {
	gotEntryResult := strings.Split(gotEntry, " ")
	wantEntryResult := strings.Split(wantEntry, " ")
	if len(gotEntryResult) != len(wantEntryResult) {
		t.Fatalf("got %d, want %d", len(gotEntryResult), len(wantEntryResult))
	}
	gotAccountNumbers := gotEntryResult[0]
	wantAccountNumbers := wantEntryResult[0]
	if gotAccountNumbers != wantAccountNumbers {
		t.Errorf("got '%s', want '%s'", gotAccountNumbers, wantAccountNumbers)
	}

	if len(gotEntryResult) > 1 {
		assertAccountNumbersStatus(t, gotEntryResult[1], wantEntryResult[1])
	}

	if len(gotEntryResult) > 2 {
		assertPossibleAccountNumbers(t, gotEntryResult[2:], wantEntryResult[2:])
	}
}

func assertPossibleAccountNumbers(t *testing.T, got, want []string) {
	t.Helper()
	gotPossibleAccountNumbers := make([]string, 0)
	for _, numbers := range got {
		numbers = strings.Trim(numbers, "[],")
		gotPossibleAccountNumbers = append(gotPossibleAccountNumbers, numbers)
	}
	wantPossibleAccountNumbers := make([]string, 0)
	for _, numbers := range want {
		numbers = strings.Trim(numbers, "[],")
		wantPossibleAccountNumbers = append(wantPossibleAccountNumbers, numbers)
	}
	sort.Strings(gotPossibleAccountNumbers)
	sort.Strings(wantPossibleAccountNumbers)
	if !reflect.DeepEqual(gotPossibleAccountNumbers, wantPossibleAccountNumbers) {
		t.Errorf("got %v, want %v", gotPossibleAccountNumbers, wantPossibleAccountNumbers)
	}
}

func assertAccountNumbersStatus(t *testing.T, got, want string) {
	t.Helper()
	assertString(t, got, want)
}

func assertString(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got '%s', want '%s'", got, want)
	}
}

func assertSliceEqual(t *testing.T, got, want []int) {
	t.Helper()
	sort.Ints(got)
	sort.Ints(want)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
