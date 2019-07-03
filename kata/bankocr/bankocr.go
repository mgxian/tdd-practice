package bankocr

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
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

func addEntryNumber(entryNumber map[string]int, entry string, number int) {
	entry = strings.ReplaceAll(entry, "\n", "")
	entryNumber[entry] = number
}

func generateEntryNumber() map[string]int {
	entryNumber := make(map[string]int, 0)
	addEntryNumber(entryNumber, ENTRY0, 0)
	addEntryNumber(entryNumber, ENTRY1, 1)
	addEntryNumber(entryNumber, ENTRY2, 2)
	addEntryNumber(entryNumber, ENTRY3, 3)
	addEntryNumber(entryNumber, ENTRY4, 4)
	addEntryNumber(entryNumber, ENTRY5, 5)
	addEntryNumber(entryNumber, ENTRY6, 6)
	addEntryNumber(entryNumber, ENTRY7, 7)
	addEntryNumber(entryNumber, ENTRY8, 8)
	addEntryNumber(entryNumber, ENTRY9, 9)

	return entryNumber
}

func parseEntry(entry string) int {
	entryNumber := generateEntryNumber()
	entry = strings.ReplaceAll(entry, "\n", "")
	if number, ok := entryNumber[entry]; ok {
		return number
	}
	return -1
}

func detectPipe(entry string) []int {
	result := make([]int, 0)
	pipePostion := []int{3, 5, 6, 8}
	for _, pos := range pipePostion {
		entryByte := []byte(entry)
		if entryByte[pos] == ' ' {
			entryByte[pos] = '|'
			number := parseEntry(string(entryByte))
			result = append(result, number)
		}
		if entryByte[pos] == '|' {
			entryByte[pos] = ' '
			number := parseEntry(string(entryByte))
			result = append(result, number)
		}
	}
	return result
}

func detectUnderscope(entry string) []int {
	result := make([]int, 0)
	underscopePostion := []int{1, 4, 7}
	for _, pos := range underscopePostion {
		entryByte := []byte(entry)
		if entryByte[pos] == ' ' {
			entryByte[pos] = '_'
			number := parseEntry(string(entryByte))
			result = append(result, number)
		}
		if entryByte[pos] == '_' {
			entryByte[pos] = ' '
			number := parseEntry(string(entryByte))
			result = append(result, number)
		}
	}
	return result
}

func smartParseEntry(entry string) []int {
	result := make([]int, 0)
	possibleNumbers := make(map[int]bool, 0)
	entry = strings.ReplaceAll(entry, "\n", "")
	number := parseEntry(entry)
	possibleNumbers[number] = true

	for _, number := range detectPipe(entry) {
		possibleNumbers[number] = true
	}

	for _, number := range detectUnderscope(entry) {
		possibleNumbers[number] = true
	}

	for k := range possibleNumbers {
		if k != -1 {
			result = append(result, k)
		}
	}
	return result
}

func splitEntry(aStringLine string) []string {
	result := make([]string, 0)
	offset := 0
	numberLength := 3
	numberLines := strings.Split(aStringLine, "\n")
	if len(numberLines) > 3 {
		numberLines = numberLines[1:4]
	}
	numberCount := len(numberLines[0]) / numberLength
	for i := 0; i < numberCount; i++ {
		numberLine1 := numberLines[0][offset : offset+numberLength]
		numberLine2 := numberLines[1][offset : offset+numberLength]
		numberLine3 := numberLines[2][offset : offset+numberLength]
		entry := fmt.Sprintf("%s%s%s", numberLine1, numberLine2, numberLine3)
		result = append(result, entry)
		offset += numberLength
	}

	return result
}

func parseStringLine(aStringLine string) []int {
	result := make([]int, 0)
	for _, entry := range splitEntry(aStringLine) {
		number := parseEntry(entry)
		result = append(result, number)
	}
	return result
}

func readStringLine(buf *bufio.Scanner) (lines []string, isEnd bool) {
	entryLineCount := 4
	for i := 0; i < entryLineCount; i++ {
		hasMoreLine := buf.Scan()
		if !hasMoreLine {
			return
		}
		line := buf.Text()
		lines = append(lines, line)
	}
	return
}

func splitStringLine(aFilePath string) []string {
	result := make([]string, 0)
	f, err := os.Open(aFilePath)
	if err != nil {
		return result
	}
	defer f.Close()
	buf := bufio.NewScanner(f)
	for {
		lines, _ := readStringLine(buf)
		if len(lines) < 3 {
			return result
		}
		if len(lines) >= 4 {
			lines = lines[:3]
		}
		stringLine := strings.Join(lines, "\n")
		result = append(result, stringLine)
	}
}

func parseNumbersFromFile(aFilePath string) [][]int {
	result := make([][]int, 0)
	for _, aStringLine := range splitStringLine(aFilePath) {
		numbers := parseStringLine(aStringLine)
		result = append(result, numbers)
	}
	return result
}

func isValidAccountNumbers(accountNumbers []int) bool {
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

func getAccountNumbersAndStatus(numbers []int) (result, status string) {
	for _, number := range numbers {
		if number == -1 {
			status = " ILL"
			result += "?"
			continue
		}
		numberString := strconv.Itoa(number)
		result += numberString
	}
	return
}

func parseAndOutputEntry(aFilePath string, w io.Writer) {
	for _, accountNumbers := range parseNumbersFromFile(aFilePath) {
		result := ""
		accountNumbersString, status := getAccountNumbersAndStatus(accountNumbers)
		result += accountNumbersString
		if status == "" && !isValidAccountNumbers(accountNumbers) {
			status = " ERR"
		}
		result += status + "\n"
		w.Write([]byte(result))
	}
}

func isCorrectedAccountNumbers(numbers []int) bool {
	for _, number := range numbers {
		if number == -1 {
			return false
		}
	}
	return true
}

func possibleAccountNumbers(entry string, i int, accountNumbers []int) (result [][]int) {
	for _, number := range smartParseEntry(entry) {
		numbers := make([]int, len(accountNumbers))
		copy(numbers, accountNumbers)
		numbers[i] = number
		if isCorrectedAccountNumbers(numbers) {
			result = append(result, numbers)
		}
	}
	return
}

func smartParseStringLine(aStringLine string) [][]int {
	result := make([][]int, 0)
	accountNumbers := parseStringLine(aStringLine)

	for i, entry := range splitEntry(aStringLine) {
		for _, numbers := range possibleAccountNumbers(entry, i, accountNumbers) {
			result = append(result, numbers)
		}
	}
	return result
}

func numbersString(numbers []int) string {
	result := ""
	for _, number := range numbers {
		if number == -1 {
			result += "?"
			continue
		}
		result += strconv.Itoa(number)
	}
	return result
}

func smartParseAndOutputEntry(aFilePath string, w io.Writer) {
	for _, line := range splitStringLine(aFilePath) {
		numbers := parseStringLine(line)
		originAccountNumbers := numbersString(numbers)
		if !strings.Contains(originAccountNumbers, "?") && isValidAccountNumbers(numbers) {
			result := originAccountNumbers + "\n"
			w.Write([]byte(result))
			continue
		}
		guessAccountNumbers(line, w)
	}
}

func possibleValidAccountNumbers(line string) []string {
	validAccountNumbers := make([]string, 0)
	for _, numbers := range smartParseStringLine(line) {
		aNumbersString := numbersString(numbers)
		if isValidAccountNumbers(numbers) {
			validAccountNumbers = append(validAccountNumbers, aNumbersString)
		}
	}
	return validAccountNumbers
}

func guessAccountNumbers(line string, w io.Writer) {
	originAccountNumbers := numbersString(parseStringLine(line))
	possibleAccountNumbers := possibleValidAccountNumbers(line)

	if len(possibleAccountNumbers) == 0 {
		status := "ILL"
		result := originAccountNumbers + " " + status + "\n"
		w.Write([]byte(result))
		return
	}

	if len(possibleAccountNumbers) == 1 {
		result := possibleAccountNumbers[0] + "\n"
		w.Write([]byte(result))
		return
	}

	status := "AMB"
	result := originAccountNumbers + " " + status
	result += " ["
	for _, aNumbersString := range possibleAccountNumbers {
		result += fmt.Sprintf("'%s', ", aNumbersString)
	}
	result = result[:len(result)-2]
	result += "]\n"
	w.Write([]byte(result))
}
