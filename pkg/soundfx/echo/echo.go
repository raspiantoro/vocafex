package echo

import (
	"encoding/binary"
	"time"

	"github.com/raspiantoro/vocafex/pkg/audio/processor"
)

type EchoFx struct {
	sampleRate float64
	delay      time.Duration
	order      binary.ByteOrder
	iteration  int
	buffer     []float32
}

func NewEcho(delay time.Duration, sampleRate float64, order binary.ByteOrder) *EchoFx {
	return &EchoFx{
		sampleRate: sampleRate,
		delay:      delay,
		order:      order,
		buffer:     make([]float32, int(sampleRate*delay.Seconds())),
	}
}

func (e *EchoFx) Process(next processor.SoundProcessor) processor.SoundProcessor {
	return processor.ProcessFunc(func(buffer *processor.SoundBuffer) {
		out := make([]float32, len(buffer.In))

		for i := range out {
			out[i] = (buffer.In[i] * .7) + (e.buffer[e.iteration] * .7)
			e.buffer[e.iteration] = buffer.In[i]
			e.iteration = (e.iteration + 1) % len(e.buffer)
		}

		buffer.In = out

		next.Process(buffer)
	})
}

// func (e *EchoFx) Process(ctx context.Context, in []float32) (out []float32) {

// 	out = make([]float32, len(e.buffer))

// 	for i := range out {
// 		out[i] = .7 * e.buffer[e.iteration]
// 		e.buffer[e.iteration] = in[i]
// 		e.iteration = (e.iteration + 1) % len(e.buffer)
// 	}

// 	return
// }
