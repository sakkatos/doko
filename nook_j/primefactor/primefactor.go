package primefactor

func PrimeFactor(number int) []int {
	if number == 1{
		return []int{}
	}
		
	for i := 2;; i++ {
		if number % i == 0 {
			return append([]int{i}, PrimeFactor(number/i)...) 
		}
	}
}
