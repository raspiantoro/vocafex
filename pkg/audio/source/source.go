package source

import (
	"bytes"
	"context"
	"fmt"
)

type audioInput interface {
	init() error
	start() error
	useDefaultConfig()
	capture(ctx context.Context) (chunk chan bytes.Buffer)
}

type AudioSource struct {
	inputType AudioInputType
	input     audioInput
}

func NewAudioSource(inputType AudioInputType, opts ...Options) (source *AudioSource, err error) {
	source = &AudioSource{
		inputType: inputType,
	}

	for _, opt := range opts {
		err = opt(source)
		if err != nil {
			return
		}
	}

	if source.input == nil {
		input, err := getAudioInput(inputType)
		if err != nil {
			return nil, err
		}
		input.useDefaultConfig()
		source.input = input
	}

	err = source.input.init()

	return
}

func getAudioInput(inputType AudioInputType) (input audioInput, err error) {
	switch inputType {
	case AudioInputTypeMic:
		return new(micInput), nil
	case AudioInputTypeWav:
		fallthrough
	case AudioInputTypeMp3:
		fallthrough
	default:
		return nil, fmt.Errorf("%s not yet implemented", inputType)
	}
}

func (a *AudioSource) Start() error {
	return a.input.start()
}

func (a *AudioSource) Capture(ctx context.Context) (chunk chan bytes.Buffer) {
	return a.input.capture(ctx)
}
