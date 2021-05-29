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
	Buffer     interface{}
}

type micInput struct {
	numChannel int
	sampleRate float64
	frameSize  int
	buffer     interface{}
	stream     *portaudio.Stream
}

func (mi *micInput) init() (err error) {
	mi.stream, err = portaudio.OpenDefaultStream(mi.numChannel, 0, mi.sampleRate, mi.frameSize, mi.buffer)
	return
}

func (mi *micInput) useDefaultConfig() {
	buffer := make([]int16, 8196)
	mi.numChannel = 1
	mi.sampleRate = 16000
	mi.frameSize = len(buffer)
	mi.buffer = buffer
}

func (mi *micInput) start() error {
	return mi.stream.Start()
}

func (mi *micInput) capture(ctx context.Context) (chunk chan bytes.Buffer) {

loopCapture:
	for {
		select {
		case <-ctx.Done():
			break loopCapture
		default:
			buff, err := mi.readBuffer()
			if err != nil {
				log.Println(err)
			}
			chunk <- buff
		}
	}

	return
}

func (mi *micInput) readBuffer() (buff bytes.Buffer, err error) {
	buff.Reset()

	fmt.Println("Read Buffer")

	err = mi.stream.Read()
	if err != nil {
		return
	}

	err = binary.Write(&buff, binary.LittleEndian, mi.buffer)

	return
}
