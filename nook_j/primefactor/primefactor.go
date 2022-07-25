package primefactor

func PrimeFactor(number int) []int {
	if number < 0{
		for i := -2;; i-- {
			if number % i == 0 {
				return append([]int{i}, PrimeFactor(number/i)...) 
			}
		}
	}
	
	if number == 1 || number == 0 || number == -1{
		return []int{}
	}
		
	for i := 2;; i++ {
		if number % i == 0 {
			return append([]int{i}, PrimeFactor(number/i)...) 
		}
	}
}
