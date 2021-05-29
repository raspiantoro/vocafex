package source

import (
	"bytes"
	"context"
)

type audioInput interface {
	init() error
	start() error
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
		// use mic as default input
		input := new(micInput)
		input.useDefaultConfig()
		source.input = input
	}

	err = source.input.init()

	return
}

func (a *AudioSource) Start() error {
	return a.input.start()
}

func (a *AudioSource) Capture(ctx context.Context) (chunk chan bytes.Buffer) {
	return a.input.capture(ctx)
}
