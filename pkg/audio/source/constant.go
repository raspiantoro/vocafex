package source

type AudioInputType int

const (
	AudioInputTypeMic AudioInputType = iota + 1
	AudioInputTypeWav
	AudioInputTypeMp3
)

func (a AudioInputType) String() string {
	switch a {
	case AudioInputTypeMic:
		return "mic input type"
	case AudioInputTypeWav:
		return "wav input type"
	case AudioInputTypeMp3:
		return "mp3 input type"
	default:
		return "invalid input type"
	}
}
