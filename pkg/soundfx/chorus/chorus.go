package chorus

import (
	"encoding/binary"
	"time"

	"github.com/raspiantoro/vocafex/pkg/audio/processor"
	"github.com/raspiantoro/vocafex/pkg/dsp/window"
)

type Chorus struct {
	sampleRate float64
	delay      time.Duration
	order      binary.ByteOrder
	iteration  int
	buffer     []float32
}

func NewChorus(delay time.Duration, sampleRate float64, order binary.ByteOrder) *Chorus {
	return &Chorus{
		sampleRate: sampleRate,
		delay:      delay,
		order:      order,
		buffer:     make([]float32, int(sampleRate*delay.Seconds())),
	}
}

func (c *Chorus) Process(next processor.SoundProcessor) processor.SoundProcessor {
	return processor.ProcessFunc(func(buffer *processor.SoundBuffer) {
		out := make([]float32, len(buffer.In))

		delaySample := window.Bartlett(len(out))

		for i := range out {
			out[i] = (buffer.In[i]) + buffer.In[i-int(delaySample[i])]
		}

		buffer.In = out

		next.Process(buffer)
	})
}
