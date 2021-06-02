package processor

import (
	"context"
)

type SoundBuffer struct {
	In  []float32
	Out []float32
}

type ProcessFunc func(buffer *SoundBuffer)

func (p ProcessFunc) Process(buffer *SoundBuffer) {
	p(buffer)
}

type SoundProcessor interface {
	Process(buffer *SoundBuffer)
}

type AudioProcessor struct {
	buffer    *SoundBuffer
	processFn []func(next SoundProcessor) SoundProcessor
}

func NewAudioProcessor() *AudioProcessor {
	b := &SoundBuffer{
		In:  []float32{},
		Out: []float32{},
	}

	a := &AudioProcessor{
		buffer:    b,
		processFn: []func(next SoundProcessor) SoundProcessor{},
	}

	a.processFn = append(a.processFn, a.defaultProcessor)

	return a
}

func (a *AudioProcessor) Register(processors ...func(SoundProcessor) SoundProcessor) {
	a.processFn = append(a.processFn, processors...)
}

func (a *AudioProcessor) Stream(ctx context.Context, in <-chan []float32, out chan<- []float32) {
	for {
		select {
		case <-ctx.Done():
			return
		case buffIn := <-in:
			a.buffer.In = buffIn

			a.Process(a.buffer)
			out <- a.buffer.Out
		}
	}
}

func (a *AudioProcessor) Process(buffer *SoundBuffer) {
	var sp SoundProcessor
	sp = a

	for _, p := range a.processFn {
		sp = p(sp)
	}

	sp.Process(buffer)
}

func (a *AudioProcessor) defaultProcessor(s SoundProcessor) SoundProcessor {
	return ProcessFunc(func(buffer *SoundBuffer) {
		buffer.Out = make([]float32, len(buffer.In))
		for i := range buffer.In {
			buffer.Out[i] = buffer.In[i]
		}
	})
}
