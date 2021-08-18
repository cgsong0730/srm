package w1

import "fmt"

func GeneratePrimeNumber(goal int) error {
	count := 0
	sum := 0

	sum = sum + goal
	fmt.Print("Prime Numbers: ")
	for i := 1; i <= goal; i++ {
		count = 0
		for j := 1; j <= i; j++ {
			if i%j == 0 {
				count += 1
			}
		}
		if count == 2 {
			fmt.Print(i, " ")
		}
	}
	fmt.Print("\n")
	return nil
}
