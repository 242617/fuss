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
	sampleRate              float64
	left, leftPhase         float64
	right, rightPhase       float64
	leftVolume, rightVolume int
	stream                  *portaudio.Stream
}

func (s *StereoSine) SetLeft(left int) {
	s.left = float64(left) / s.sampleRate
}

func (s *StereoSine) SetRight(right int) {
	s.right = float64(right) / s.sampleRate
}

func (s *StereoSine) SetVolume(volume int) {
	s.leftVolume, s.rightVolume = volume, volume
}

func (s *StereoSine) Play() error {
	return s.stream.Start()
}

func (s *StereoSine) Stop() error {
	return s.stream.Stop()
}

func (s *StereoSine) process(out [][]float32) {
	for i := range out[0] {
		out[0][i] = float32(math.Sin(2*math.Pi*s.leftPhase)) * float32(s.leftVolume) / 100
		_, s.leftPhase = math.Modf(s.leftPhase + s.left)
		out[1][i] = float32(math.Sin(2*math.Pi*s.rightPhase)) * float32(s.rightVolume) / 100
		_, s.rightPhase = math.Modf(s.rightPhase + s.right)
	}
}

func cycle(min, max int) func() int {
	if min > max {
		min, max = max, min
	}
	cur := min
	inc := true

	return func() int {
		if inc {
			cur += 1
			if cur == max {
				inc = false
			}
		} else {
			cur -= 1
			if cur == min {
				inc = true
			}
		}
		return cur
	}
}
