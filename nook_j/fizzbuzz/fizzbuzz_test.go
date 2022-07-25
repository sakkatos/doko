package fizzbuzz

import (
	"testing"
	// "github.com/stretchr/testify/assert"
)

func TestFizzBUZZ_1_Return_1(t *testing.T) {
	answer := FizzBuzz(1)
	expected := "1"
	// assert.Equal(t, answer, expected, "The two words should be the same.")
	if answer != expected {
		t.Errorf("got %s, wanted %s", answer, expected)
	}
}

func TestFizzBUZZ_2_Return_2(t *testing.T) {
	answer := FizzBuzz(2)
	expected := "2"

	if answer != expected {
		t.Errorf("got %s, wanted %s", answer, expected)
	}
}

func TestFizzBUZZ_3_Return_fizz(t *testing.T) {
	answer := FizzBuzz(3)
	expected := "fizz"

	if answer != expected {
		t.Errorf("got %s, wanted %s", answer, expected)
	}
}

func TestFizzBUZZ_5_Return_buzz(t *testing.T) {
	answer := FizzBuzz(5)
	expected := "buzz"

	if answer != expected {
		t.Errorf("got %s, wanted %s", answer, expected)
	}
}

func TestFizzBUZZ_6_Return_fizz(t *testing.T) {
	answer := FizzBuzz(6)
	expected := "fizz"

	if answer != expected {
		t.Errorf("got %s, wanted %s", answer, expected)
	}
}

func TestFizzBUZZ_10_Return_buzz(t *testing.T) {
	answer := FizzBuzz(10)
	expected := "buzz"

	if answer != expected {
		t.Errorf("got %s, wanted %s", answer, expected)
	}
}

func TestFizzBUZZ_15_Return_fizzbuzz(t *testing.T) {
	answer := FizzBuzz(15)
	expected := "fizzbuzz"

	if answer != expected {
		t.Errorf("got %s, wanted %s", answer, expected)
	}
}

func TestFizzBUZZ_30_Return_fizzbuzz(t *testing.T) {
	answer := FizzBuzz(30)
	expected := "fizzbuzz"

	if answer != expected {
		t.Errorf("got %s, wanted %s", answer, expected)
	}
}
