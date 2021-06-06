package filter

import (
	"math"

	"github.com/raspiantoro/vocafex/pkg/audio/processor"
)

type HighpassFilterConfig struct {
	Cutoff    float32
	CutoffMod float32
	Resonance float32
	Gain      float32
}

type HighpassFilter struct {
	cutoff                             float32
	cutoffMod                          float32
	resonance                          float32
	feedback                           float32
	gain                               float32
	sample0, sample1, sample2, sample3 float32
}

func NewHighpassFilter(cfg HighpassFilterConfig) *HighpassFilter {
	return &HighpassFilter{
		gain:      cfg.Gain,
		cutoff:    cfg.Cutoff,
		cutoffMod: cfg.CutoffMod,
		feedback:  0,
		resonance: cfg.Resonance,
		sample0:   0,
		sample1:   0,
		sample2:   0,
		sample3:   0,
	}
}

func (hp *HighpassFilter) Process(next processor.SoundProcessor) processor.SoundProcessor {
	return processor.ProcessFunc(func(buffer *processor.SoundBuffer) {
		out := make([]float32, len(buffer.In))
		cutoff := hp.calculatedCutoff()
		hp.calculateFeedback()

		for i := range out {
			hp.sample0 += cutoff * (buffer.In[i] - hp.sample0 + hp.feedback*(hp.sample0-hp.sample1))
			hp.sample1 += cutoff * (hp.sample0 - hp.sample1)
			hp.sample2 += cutoff * (hp.sample1 - hp.sample2)
			hp.sample3 += cutoff * (hp.sample2 - hp.sample3)

			out[i] = (buffer.In[i] - hp.sample3) * hp.gain
		}

		buffer.In = out

		next.Process(buffer)
	})
}

func (hp *HighpassFilter) calculatedCutoff() float32 {
	return float32(math.Max(math.Min(float64(hp.cutoff)+float64(hp.cutoffMod), 0.99), 0.01))
}

func (hp *HighpassFilter) calculateFeedback() {
	hp.feedback = float32(math.Min(float64(hp.resonance+hp.resonance/(1.0-hp.calculatedCutoff())), 2.0))
}
