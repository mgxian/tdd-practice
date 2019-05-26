package fizzbuzz

import "strconv"

// FizzBuzz is a number game
func FizzBuzz(number int) string {
	if number%3 == 0 && number%5 == 0 {
		return "fizzbuzz"
	}

	if number%3 == 0 {
		return "fizz"
	}

	if number%5 == 0 {
		return "buzz"
	}

	return strconv.Itoa(number)
}
