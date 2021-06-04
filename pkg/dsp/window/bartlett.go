package window

import "fmt"

func Bartlett(len int) []int32 {

	output := make([]int32, len)

	N := len - 1

	coefficients := 2 / N
	n := 0

	fmt.Println(N)

	for ; n <= N/2; n++ {
		output[n] = int32(n * coefficients)
	}

	for ; n <= N; N++ {
		output[n] = int32(2 - n*int(coefficients))
	}

	return output
}
