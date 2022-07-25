package primefactor

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func Test_PrimefactorOf2_shouldbe_arrayof2(t *testing.T) {
	answer := PrimeFactor(2)
	expected := []int{2}
	assert.Equal(t, answer, expected)
}

func Test_PrimefactorOf3_shouldbe_arrayof3(t *testing.T) {
	answer := PrimeFactor(3)
	expected := []int{3}
	assert.Equal(t, answer, expected)
}

func Test_PrimefactorOf9_shouldbe_arrayof3and3(t *testing.T) {
	answer := PrimeFactor(9)
	expected := []int{3, 3}
	assert.Equal(t, answer, expected, "คำตอบต้องเป็น %v แต่ได้ %v", expected, answer)
}
func Test_PrimefactorOf6_shouldbe_arrayof2and3(t *testing.T) {
	answer := PrimeFactor(6)
	expected := []int{2, 3}
	assert.Equal(t, answer, expected, "คำตอบต้องเป็น %v แต่ได้ %v", expected, answer)
}

func Test_PrimefactorOf4_shouldbe_arrayof2and2(t *testing.T) {
	answer := PrimeFactor(4)
	expected := []int{2, 2}
	assert.Equal(t, answer, expected, "คำตอบต้องเป็น %v แต่ได้ %v", expected, answer)
}

func Test_PrimefactorOf25_shouldbe_arrayof5and5(t *testing.T) {
	answer := PrimeFactor(25)
	expected := []int{5, 5}
	assert.Equal(t, answer, expected, "คำตอบต้องเป็น %v แต่ได้ %v", expected, answer)
}

func Test_PrimefactorOf27_shouldbe_arrayof3and3and3(t *testing.T) {
	answer := PrimeFactor(27)
	expected := []int{3,3,3}
	assert.Equal(t, answer, expected, "คำตอบต้องเป็น %v แต่ได้ %v", expected, answer)
}

func Test_PrimefactorOf12_shouldbe_arrayof2and2and3(t *testing.T) {
	answer := PrimeFactor(12)
	expected := []int{2,2,3}
	assert.Equal(t, answer, expected, "คำตอบต้องเป็น %v แต่ได้ %v", expected, answer)
}
func Test_PrimefactorOf1_shouldbe_emptyArray(t *testing.T) {
	answer := PrimeFactor(1)
	expected := []int{}
	assert.Equal(t, answer, expected, "คำตอบต้องเป็น %v แต่ได้ %v", expected, answer)
}

func Test_PrimefactorOfminus9_shouldbe_minus3and3(t *testing.T) {
	answer := PrimeFactor(-9)
	expected := []int{-3,3}
	assert.Equal(t, answer, expected, "คำตอบต้องเป็น %v แต่ได้ %v", expected, answer)
}

