package sink

import (
	"bytes"
	"encoding/binary"
	"log"

	"github.com/gordonklaus/portaudio"
)

type SpeakerConfig struct {
	NumChannel int
	SampleRate float64
	FrameSize  int
	Order      binary.ByteOrder
	Buffer     interface{}
}

type speakerOutput struct {
	cfg    SpeakerConfig
	stream *portaudio.Stream
}

func (s *speakerOutput) init() (err error) {
	s.stream, err = portaudio.OpenDefaultStream(0, s.cfg.NumChannel, s.cfg.SampleRate, s.cfg.FrameSize, s.cfg.Buffer)
	return
}

func (s *speakerOutput) useDefaultConfig() {
	buffer := make([]int16, 8196)
	s.cfg = SpeakerConfig{}
	s.cfg.NumChannel = 1
	s.cfg.SampleRate = 16000
	s.cfg.FrameSize = len(buffer)
	s.cfg.Order = binary.LittleEndian
	s.cfg.Buffer = buffer
}

func (s *speakerOutput) start() error {
	return s.stream.Start()
}

func (s *speakerOutput) receive(chunk <-chan bytes.Buffer) (err error) {
	for buff := range chunk {
		err = binary.Read(&buff, s.cfg.Order, s.cfg.Buffer)
		if err != nil {
			log.Println(err)
		}

		err = s.stream.Write()
		if err != nil {
			log.Println(err)
		}
	}

	return
}
