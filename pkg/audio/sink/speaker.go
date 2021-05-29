package sink

import (
	"context"
	"encoding/binary"

	"github.com/gordonklaus/portaudio"
)

type SpeakerConfig struct {
	NumChannel int
	SampleRate float64
	Order      binary.ByteOrder
	Buffer     []float32
}

type speakerOutput struct {
	cfg    SpeakerConfig
	stream *portaudio.Stream
}

func (s *speakerOutput) useDefaultConfig() {
	buffer := make([]float32, 8196)
	s.cfg = SpeakerConfig{}
	s.cfg.NumChannel = 1
	s.cfg.SampleRate = 16000
	s.cfg.Order = binary.LittleEndian
	s.cfg.Buffer = buffer
}

func (s *speakerOutput) start() (err error) {
	h, err := portaudio.DefaultHostApi()
	if err != nil {
		return
	}

	p := portaudio.LowLatencyParameters(nil, h.DefaultOutputDevice)
	p.Output.Channels = s.cfg.NumChannel
	p.SampleRate = s.cfg.SampleRate
	// p.FramesPerBuffer = 64
	sample := 0

	s.stream, err = portaudio.OpenStream(p, func(out []float32) {
		// fmt.Println(len(s.cfg.Buffer))
		for i := range out {
			// 	// fmt.Println(out[i])
			// 	// if out[i] > 0 {
			out[i] = s.cfg.Buffer[(sample+i)%len(s.cfg.Buffer)]
			// 	// }

			// out[i] = s.cfg.Buffer[i]
		}
		sample += len(out)
	})
	if err != nil {
		return
	}

	return s.stream.Start()
}

func (s *speakerOutput) receive(ctx context.Context, chunk <-chan []float32) (err error) {
	for buff := range chunk {
		select {
		case <-ctx.Done():
			s.stream.Stop()
			s.stream.Close()
			return
		default:
			// err = binary.Read(&buff, s.cfg.Order, s.cfg.Buffer)
			// if err != nil {
			// 	// log.Println(err)
			// }

			// fmt.Println(buff)

			s.cfg.Buffer = make([]float32, len(buff))

			s.cfg.Buffer = buff
			err = s.stream.Write()
			if err != nil {
				// log.Println(err)
			}
		}

	}

	return
}
