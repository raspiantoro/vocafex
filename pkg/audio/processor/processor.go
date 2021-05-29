package processor

type SoundProcessor interface {
	ProcessAudio(in []float32)
	GetBuffer() []float32
}

type AudioProcessor struct {
	processor SoundProcessor
}

func NewAudioProcessor() *AudioProcessor {
	return new(AudioProcessor)
}

func (a *AudioProcessor) Register(processor SoundProcessor) {
	a.processor = processor
}

func (a *AudioProcessor) Process(in []float32) {
	a.processor.ProcessAudio(in)
}

func (a *AudioProcessor) GetProcessor() SoundProcessor {
	return a.processor
}

func (a *AudioProcessor) GetBuffer() (buff []float32) {
	return a.processor.GetBuffer()
}
