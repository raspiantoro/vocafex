package sink

import (
	"bytes"
	"fmt"
)

type audioOutput interface {
	init() (err error)
	start() (err error)
	useDefaultConfig()
	receive(chunk <-chan bytes.Buffer) error
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

	err = sink.output.init()

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

func (a *AudioSink) Receive(chunk chan bytes.Buffer) error {
	return a.output.receive(chunk)
}
