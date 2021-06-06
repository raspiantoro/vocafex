package timebased

import (
	"time"

	"github.com/raspiantoro/vocafex/pkg/audio/processor"
	"github.com/raspiantoro/vocafex/pkg/buffer"
)

type Delay struct {
	sampleRate     float64
	delay          time.Duration
	delayRate      int
	gain           float32
	samplePerBlock int
	buffer         *buffer.CircularBuffer
}

func NewDelay(delay time.Duration, delayRate int, gain float32, sampleRate float64, samplePerBlock int) *Delay {
	return &Delay{
		sampleRate:     sampleRate,
		delay:          delay,
		delayRate:      delayRate,
		gain:           gain,
		samplePerBlock: samplePerBlock,
		buffer:         buffer.NewCircularBuffer(int(delay.Seconds())*(int(sampleRate)+samplePerBlock), delayRate),
	}
}

func (d *Delay) Process(next processor.SoundProcessor) processor.SoundProcessor {
	return processor.ProcessFunc(func(buffer *processor.SoundBuffer) {
		out := make([]float32, len(buffer.In))

		for i := range out {
			delaySample, _ := d.buffer.Read()
			out[i] = (buffer.In[i]) + (delaySample * d.gain)
			d.buffer.Write(buffer.In[i])
		}

		buffer.In = out

		next.Process(buffer)
	})
}
