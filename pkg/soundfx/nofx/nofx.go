package nofx

type NoFx struct{}

func (n *NoFx) ProcessAudio(out chan []float32) func(in []float32) {
	return func(in []float32) {
		newBuff := make([]float32, len(in))
		sample := 0
		for i := range in {
			newBuff[i] = in[(sample+i)%len(in)]
		}

		out <- newBuff
	}
}
