package source

import (
	"context"
	"encoding/binary"
	"log"

	"github.com/gordonklaus/portaudio"
)

type MicConfig struct {
	NumChannel int
	SampleRate float64
	Order      binary.ByteOrder
	Buffer     []float32
}

type micInput struct {
	cfg      MicConfig
	streamer *portaudio.Stream
}

func (m *micInput) useDefaultConfig() {
	buffer := make([]float32, 8196)
	m.cfg = MicConfig{}
	m.cfg.NumChannel = 1
	m.cfg.SampleRate = 16000
	m.cfg.Order = binary.LittleEndian
	m.cfg.Buffer = buffer
}

func (m *micInput) start() (err error) {
	h, err := portaudio.DefaultHostApi()
	if err != nil {
		return
	}

	p := portaudio.LowLatencyParameters(h.DefaultInputDevice, nil)
	p.Input.Channels = m.cfg.NumChannel
	p.SampleRate = m.cfg.SampleRate
	p.FramesPerBuffer = 8

	m.streamer, err = portaudio.OpenStream(p, m.cfg.Buffer)
	if err != nil {
		return
	}

	return m.streamer.Start()
}

func (m *micInput) stream(ctx context.Context, chunk chan<- []float32) {

	go func() {
		for {
			select {
			case <-ctx.Done():
				m.streamer.Close()
				m.streamer.Stop()
				return
			default:
				buff, err := m.readBuffer()
				if err != nil {
					log.Println(err)
				}

				chunk <- buff
			}
		}

	}()
}

func (m *micInput) readBuffer() (buff []float32, err error) {

	err = m.streamer.Read()
	if err != nil {
		return
	}

	buff = m.cfg.Buffer

	return
}
