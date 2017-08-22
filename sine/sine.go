package sine

import (
	"math"

	"github.com/gordonklaus/portaudio"
)

func NewStereoSine(left, right int, sampleRate int) (s *StereoSine, err error) {
	s = &StereoSine{}
	s.sampleRate = float64(sampleRate)

	if s.stream, err = portaudio.OpenDefaultStream(0, 2, float64(s.sampleRate), 0, s.process); err != nil {
		return
	}
	s.SetLeft(left)
	s.SetRight(right)

	return
}

type StereoSine struct {
	sampleRate        float64
	left, leftPhase   float64
	right, rightPhase float64
	stream            *portaudio.Stream
}

func (s *StereoSine) SetLeft(left int) {
	s.left = float64(left) / s.sampleRate
}

func (s *StereoSine) SetRight(right int) {
	s.right = float64(right) / s.sampleRate
}

func (s *StereoSine) Play() error {
	return s.stream.Start()
}

func (s *StereoSine) Stop() error {
	return s.stream.Stop()
}

func (s *StereoSine) process(out [][]float32) {
	for i := range out[0] {
		out[0][i] = float32(math.Sin(2 * math.Pi * s.leftPhase))
		_, s.leftPhase = math.Modf(s.leftPhase + s.left)
		out[1][i] = float32(math.Sin(2 * math.Pi * s.rightPhase))
		_, s.rightPhase = math.Modf(s.rightPhase + s.right)
	}
}