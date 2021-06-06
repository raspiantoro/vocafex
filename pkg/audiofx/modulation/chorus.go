package modulation

import (
	"time"

	"github.com/raspiantoro/vocafex/pkg/audio/processor"
	"github.com/raspiantoro/vocafex/pkg/buffer"
	"github.com/raspiantoro/vocafex/pkg/dsp/window"
)

type Chorus struct {
	sampleRate     float64
	delay          time.Duration
	delayRate      int
	gain           float32
	samplePerBlock int
	delayBuffer    *buffer.CircularBuffer
}

func NewChorus(delay time.Duration, delayRate int, gain float32, sampleRate float64, samplePerBlock int) *Chorus {
	return &Chorus{
		sampleRate:     sampleRate,
		delay:          delay,
		delayRate:      delayRate,
		gain:           gain,
		samplePerBlock: samplePerBlock,
		delayBuffer:    buffer.NewCircularBuffer(int(delay.Seconds())*(int(sampleRate)+samplePerBlock), delayRate),
	}
}

func (c *Chorus) Process(next processor.SoundProcessor) processor.SoundProcessor {
	return processor.ProcessFunc(func(buffer *processor.SoundBuffer) {
		out := make([]float32, len(buffer.In))

		offset := window.Bartlett(len(out))

		for i := range out {
			delaySample, _ := c.delayBuffer.Read()
			out[i] = buffer.In[i] + (delaySample * c.gain * offset[i])
			c.delayBuffer.Write(buffer.In[i])
		}

		buffer.In = out

		next.Process(buffer)
	})
}
