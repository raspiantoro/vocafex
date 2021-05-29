package sink

import "fmt"

type Options func(*AudioSink) error

func WithSpeakerConfig(cfg SpeakerConfig) Options {
	return func(a *AudioSink) (err error) {
		if a.outputType != AudioOutputTypeSpeaker {
			err = fmt.Errorf("can't use %s with speaker output config", a.outputType)
			return
		}

		a.output = &speakerOutput{
			cfg: cfg,
		}

		return
	}
}
