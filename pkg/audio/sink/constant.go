package sink

type AudioOutputType int

const (
	AudioOutputTypeSpeaker AudioOutputType = iota + 1
	AudioOutputTypeWav
	AudioOutputTypeMp3
)

func (a AudioOutputType) String() string {
	switch a {
	case AudioOutputTypeSpeaker:
		return "speaker output type"
	case AudioOutputTypeWav:
		return "wav output type"
	case AudioOutputTypeMp3:
		return "mp3 output type"
	default:
		return "invalid sink type"
	}
}
