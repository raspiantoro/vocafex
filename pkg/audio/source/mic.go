package source

import (
	"bytes"
	"context"
	"encoding/binary"
	"log"

	"github.com/gordonklaus/portaudio"
	"github.com/raspiantoro/vocafex/pkg/audio/processor"
)

type MicConfig struct {
	NumChannel int
	SampleRate float64
	Order      binary.ByteOrder
	Buffer     []float32
}

type micInput struct {
	cfg       MicConfig
	stream    *portaudio.Stream
	processor processor.SoundProcessor
}

func (m *micInput) useDefaultConfig() {
	buffer := make([]float32, 8196)
	m.cfg = MicConfig{}
	m.cfg.NumChannel = 1
	m.cfg.SampleRate = 16000
	m.cfg.Order = binary.LittleEndian
	m.cfg.Buffer = buffer
}

func (m *micInput) RegisterProcessor(processor processor.SoundProcessor) (err error) {
	m.processor = processor
	return
}

func (m *micInput) start() (err error) {
	h, err := portaudio.DefaultHostApi()
	if err != nil {
		return
	}

	p := portaudio.LowLatencyParameters(h.DefaultInputDevice, nil)
	p.Input.Channels = m.cfg.NumChannel
	p.SampleRate = m.cfg.SampleRate

	// p.FramesPerBuffer = 64

	m.stream, err = portaudio.OpenStream(p, m.cfg.Buffer)
	if err != nil {
		return
	}

	return m.stream.Start()
}

func (m *micInput) capture(ctx context.Context) (chunk chan bytes.Buffer) {

	chunk = make(chan bytes.Buffer)

	go func() {
	loopCapture:
		for {
			select {
			case <-ctx.Done():
				m.stream.Close()
				m.stream.Stop()
				break loopCapture
			default:
				_, err := m.readBuffer()
				if err != nil {
					log.Println(err)
				}
				// chunk <- buff
			}
		}

	}()

	return
}

func (m *micInput) readBuffer() (buff bytes.Buffer, err error) {
	buff = bytes.Buffer{}
	buff.Reset()

	err = m.stream.Read()
	if err != nil {
		return
	}

	// data := make(chan []float32)
	// fmt.Println(m.cfg.Buffer)
	m.processor.ProcessAudio(m.cfg.Buffer)

	// newBuffer := <-data

	// err = binary.Write(&buff, m.cfg.Order, newBuffer)

	return
}
