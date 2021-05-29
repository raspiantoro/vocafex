package pipeline

import (
	"context"
	"encoding/binary"
	"time"

	"github.com/raspiantoro/vocafex/pkg/audio/processor"
	"github.com/raspiantoro/vocafex/pkg/audio/sink"
	"github.com/raspiantoro/vocafex/pkg/audio/source"
	"github.com/raspiantoro/vocafex/pkg/soundfx/echo"
)

type Pipeline struct {
	Source    *source.AudioSource
	Processor *processor.AudioProcessor
	Sink      *sink.AudioSink
}

func (p *Pipeline) Start(ctx context.Context) (err error) {

	if p.Processor == nil {
		p.Processor = new(processor.AudioProcessor)
		// nofx := new(nofx.NoFx)
		// p.Processor.Register(nofx)

		e := echo.NewEcho(time.Second/5, 44100, binary.BigEndian)
		p.Processor.Register(e)
	}

	p.Source.RegisterProcessor(p.Processor.GetProcessor())

	err = p.Source.Start()
	if err != nil {
		return
	}

	err = p.Sink.Start()
	if err != nil {
		return
	}

	dataChan := make(chan []float32)

	go p.Sink.Receive(ctx, dataChan)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				// buff := bytes.Buffer{}
				// err = binary.Read(&buff, binary.BigEndian, p.Processor.GetBuffer())
				// if err == nil {
				dataChan <- p.Processor.GetBuffer()
				// }
			}
		}

	}()

	go p.Source.Capture(ctx)

	return
}
