package window

import "fmt"

func Bartlett(len int) []float32 {

	output := make([]float32, len)

	N := float32(len - 1)

	coefficients := 2 / N
	n := 0

	fmt.Println(N)

	for ; n <= N/2; n++ {
		output[n] = float32(n * int(coefficients))
	}

	for ; n <= N; N++ {
		output[n] = 2 - float32(n*int(coefficients))
	}

	return output
}
