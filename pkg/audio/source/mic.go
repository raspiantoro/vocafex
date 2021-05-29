package source

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"log"

	"github.com/gordonklaus/portaudio"
)

type MicConfig struct {
	NumChannel int
	SampleRate float64
	FrameSize  int
	Order      binary.ByteOrder
	Buffer     interface{}
}

type micInput struct {
	cfg    MicConfig
	stream *portaudio.Stream
}

func (m *micInput) init() (err error) {
	m.stream, err = portaudio.OpenDefaultStream(m.cfg.NumChannel, 0, m.cfg.SampleRate, m.cfg.FrameSize, m.cfg.Buffer)
	return
}

func (m *micInput) useDefaultConfig() {
	buffer := make([]int16, 8196)
	m.cfg = MicConfig{}
	m.cfg.NumChannel = 1
	m.cfg.SampleRate = 16000
	m.cfg.FrameSize = len(buffer)
	m.cfg.Order = binary.LittleEndian
	m.cfg.Buffer = buffer
}

func (m *micInput) start() error {
	return m.stream.Start()
}

func (m *micInput) capture(ctx context.Context) (chunk chan bytes.Buffer) {

	chunk = make(chan bytes.Buffer)

	go func() {
	loopCapture:
		for {
			select {
			case <-ctx.Done():
				break loopCapture
			default:
				buff, err := m.readBuffer()
				if err != nil {
					log.Println(err)
				}
				chunk <- buff
			}
		}

		return
	}()

	return
}

func (m *micInput) readBuffer() (buff bytes.Buffer, err error) {
	buff = bytes.Buffer{}
	buff.Reset()

	fmt.Println("Read Buffer")

	err = m.stream.Read()
	if err != nil {
		return
	}

	err = binary.Write(&buff, m.cfg.Order, m.cfg.Buffer)

	return
}
