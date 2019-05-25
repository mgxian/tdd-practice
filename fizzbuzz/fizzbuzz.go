package fizzbuzz

import (
	"strconv"
	"strings"
)

func isDivisibleBy(number, divisor int) bool {
	return number%divisor == 0
}

func isContains(number, subNumber int) bool {
	return strings.Contains(strconv.Itoa(number), strconv.Itoa(subNumber))
}

// FizzBuzz magic number
const (
	FizzMagciNumber int = 3
	BuzzMagicNumber int = 5
)

// FizzBuzz is a number game.
func FizzBuzz(number int) string {
	if isDivisibleBy(number, FizzMagciNumber) && isDivisibleBy(number, BuzzMagicNumber) {
		return "fizzbuzz"
	}

	if isContains(number, FizzMagciNumber) && isContains(number, BuzzMagicNumber) {
		return "fizzbuzz"
	}

	if isDivisibleBy(number, FizzMagciNumber) {
		return "fizz"
	}

	if isDivisibleBy(number, BuzzMagicNumber) {
		return "buzz"
	}

	if isContains(number, FizzMagciNumber) {
		return "fizz"
	}

	if isContains(number, BuzzMagicNumber) {
		return "buzz"
	}

	return strconv.Itoa(number)
}
