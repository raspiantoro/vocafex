package source

import "fmt"

type Options func(*AudioSource) error

func WithMicConfig(cfg MicConfig) Options {
	return func(a *AudioSource) (err error) {

		if a.inputType != AudioInputTypeMic {
			err = fmt.Errorf("can't use %s with mic input config", a.inputType)
			return
		}

		a.input = &micInput{
			cfg: cfg,
		}

		return
	}
}
