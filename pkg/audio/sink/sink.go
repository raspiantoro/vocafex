package sink

import (
	"context"
	"fmt"
)

type audioOutput interface {
	start() (err error)
	useDefaultConfig()
	receive(ctx context.Context, chunk <-chan []float32) error
}

type AudioSink struct {
	outputType AudioOutputType
	output     audioOutput
}

func NewAudioSink(outputType AudioOutputType, opts ...Options) (sink *AudioSink, err error) {
	sink = &AudioSink{
		outputType: outputType,
	}

	for _, opt := range opts {
		err = opt(sink)
		if err != nil {
			return
		}
	}

	if sink.output == nil {
		output, err := getAudioOutput(outputType)
		if err != nil {
			return nil, err
		}

		output.useDefaultConfig()
		sink.output = output
	}

	return
}

func getAudioOutput(outputType AudioOutputType) (output audioOutput, err error) {
	switch outputType {
	case AudioOutputTypeSpeaker:
		return new(speakerOutput), nil
	case AudioOutputTypeWav:
		fallthrough
	case AudioOutputTypeMp3:
		fallthrough
	default:
		return nil, fmt.Errorf("%s not yet implemented", outputType)
	}
}

func (a *AudioSink) Start() error {
	return a.output.start()
}

func (a *AudioSink) Receive(ctx context.Context, chunk chan []float32) error {
	return a.output.receive(ctx, chunk)
}
