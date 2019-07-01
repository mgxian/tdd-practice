package bankocr

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// entry string
const (
	ENTRY1 string = `
   
  |
  |
`
	ENTRY2 string = `
 _ 
 _|
|_ 
`

	ENTRY3 string = `
 _ 
 _|
 _|
`

	ENTRY4 string = `
   
|_|
  |
`

	ENTRY5 string = `
 _ 
|_ 
 _|
`

	ENTRY6 string = `
 _ 
|_ 
|_|
`

	ENTRY7 string = `
 _ 
  |
  |
`

	ENTRY8 string = `
 _ 
|_|
|_|
`

	ENTRY9 string = `
 _ 
|_|
 _|
`

	ENTRY0 string = `
 _ 
| |
|_|
`
)

func parseEntry(entry string) int {
	entryNumber := make(map[string]int, 0)
	entryNumber[ENTRY0] = 0
	entryNumber[ENTRY1] = 1
	entryNumber[ENTRY2] = 2
	entryNumber[ENTRY3] = 3
	entryNumber[ENTRY4] = 4
	entryNumber[ENTRY5] = 5
	entryNumber[ENTRY6] = 6
	entryNumber[ENTRY7] = 7
	entryNumber[ENTRY8] = 8
	entryNumber[ENTRY9] = 9

	if !strings.HasPrefix(entry, "\n") {
		entry = "\n" + entry
	}

	if number, ok := entryNumber[entry]; ok {
		return number
	}
	return -1
}

func parseStringLine(aStringLine string) []int {
	result := make([]int, 0)
	offset := 0
	numberLength := 3
	numberLines := strings.Split(aStringLine, "\n")
	if len(numberLines) > 3 {
		numberLines = numberLines[1:4]
	}
	// fmt.Println(numberLines)
	numberCount := len(numberLines[0]) / numberLength
	for i := 0; i < numberCount; i++ {
		numberLine1 := numberLines[0][offset : offset+numberLength]
		numberLine2 := numberLines[1][offset : offset+numberLength]
		numberLine3 := numberLines[2][offset : offset+numberLength]
		entry := fmt.Sprintf("%s\n%s\n%s\n", numberLine1, numberLine2, numberLine3)
		number := parseEntry(entry)
		result = append(result, number)
		offset += numberLength
	}
	return result
}

func parseNumbersFromFile(aFilePath string) [][]int {
	result := make([][]int, 0)
	f, err := os.Open(aFilePath)
	if err != nil {
		return result
	}
	defer f.Close()
	entryLineCount := 4
	buf := bufio.NewScanner(f)
	done := false
	for {
		lines := make([]string, 0)
		for i := 0; i < entryLineCount; i++ {
			hasLine := buf.Scan()
			line := buf.Text()
			err := buf.Err()
			if err != nil {
				return result
			}
			if !hasLine {
				done = true
			}
			lines = append(lines, line)
		}
		if len(lines) < 3 {
			continue
		}
		if len(lines) == 4 {
			lines = lines[:3]
		}
		stringLine := strings.Join(lines, "\n")
		result = append(result, parseStringLine(stringLine))
		if done {
			return result
		}
	}
}

func validAccountNumbers(accountNumbers []int) bool {
	sum := 0
	length := len(accountNumbers)
	for i, number := range accountNumbers {
		sum += (length - i) * number
	}

	if sum%11 == 0 {
		return true
	}
	return false
}
