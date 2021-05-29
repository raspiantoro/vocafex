package echo

import (
	"encoding/binary"
	"time"
)

type EchoFx struct {
	sampleRate float64
	delay      time.Duration
	order      binary.ByteOrder
	iteration  int
	buffer     []float32
	buffChunk  []float32
}

func NewEcho(delay time.Duration, sampleRate float64, order binary.ByteOrder) *EchoFx {
	return &EchoFx{
		sampleRate: sampleRate,
		delay:      delay,
		order:      order,
		buffer:     make([]float32, int(sampleRate*delay.Seconds())),
		buffChunk:  make([]float32, int(sampleRate*delay.Seconds())),
	}
}

func (e *EchoFx) ProcessAudio(in []float32) {

	newBuff := make([]float32, len(e.buffer))
	// sample := 0
	for i := range newBuff {
		// b := in[(sample+i)%len(in)]
		newBuff[i] = .7 * e.buffer[e.iteration]
		// newBuff[i] = .7 * in[e.iteration]
		e.buffer[e.iteration] = in[i]
		e.iteration = (e.iteration + 1) % len(e.buffer)
	}

	e.buffChunk = newBuff

}

func (e *EchoFx) GetBuffer() []float32 {
	// fmt.Println(len(e.buffChunk))
	return e.buffer
}
