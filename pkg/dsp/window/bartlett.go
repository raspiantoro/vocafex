package window

import "fmt"

func Bartlett(len int) []float32 {

	output := make([]float32, len)

	N := len - 1

	coefficients := 2 / float32(N)
	fmt.Printf("%f\n", coefficients)
	n := 0

	for ; n <= N/2; n++ {
		output[n] = float32(n) * coefficients
	}

	for ; n <= N; n++ {
		output[n] = 2 - (float32(n) * coefficients)
	}

	return output
}
