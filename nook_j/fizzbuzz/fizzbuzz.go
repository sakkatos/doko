package fizzbuzz

import "strconv"

func FizzBuzz(input int) string {
	if isBuzz(input) && isFizz(input) {
		return "fizzbuzz"
	}
	if isBuzz(input) {
		return "buzz"
	}
	if isFizz(input) {
		return "fizz"
	}

	return strconv.Itoa(input)
}

func isBuzz(input int) bool {
	return input%5 == 0
}

func isFizz(input int) bool {
	return input%3 == 0
}
