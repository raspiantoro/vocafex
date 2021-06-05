package pipeline

import (
	"context"

	"github.com/raspiantoro/vocafex/pkg/audio/processor"
	"github.com/raspiantoro/vocafex/pkg/audio/sink"
	"github.com/raspiantoro/vocafex/pkg/audio/source"
	"github.com/raspiantoro/vocafex/pkg/audioprocessor/filter"
)

type Pipeline struct {
	Source    *source.AudioSource
	Processor *processor.AudioProcessor
	Sink      *sink.AudioSink
}

func (p *Pipeline) Start(ctx context.Context) (err error) {
	sourceStreamData := make(chan []float32)
	sinkStreamData := make(chan []float32)

	p.Processor = processor.NewAudioProcessor()
	// nofx := new(nofx.NoFx)
	// p.Processor.Register(nofx)

	// c := chorus.NewChorus(time.Second/3, 44100, binary.BigEndian)
	// p.Processor.Register(c.Process)

	// d := fx.NewDelay(time.Second*2, 10000, .7, 44100, 512)
	// p.Processor.Register(d.Process)

	// d2 := fx.NewDelay(time.Second*5, 10000, .5, 44100, 512)
	// p.Processor.Register(d2.Process)

	hfConfig := filter.HighpassFilterConfig{
		Cutoff:    .6,
		CutoffMod: .2,
		Resonance: 6,
		Gain:      80,
	}

	hf := filter.NewHighpassFilter(hfConfig)
	p.Processor.Register(hf.Process)

	err = p.Source.Start()
	if err != nil {
		return
	}

	err = p.Sink.Start()
	if err != nil {
		return
	}

	go p.Source.Stream(ctx, sourceStreamData)

	go p.Processor.Stream(ctx, sourceStreamData, sinkStreamData)

	go p.Sink.Receive(ctx, sinkStreamData)

	return
}
