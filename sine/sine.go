package sine

import (
	"math"

	"github.com/gordonklaus/portaudio"
)

func NewStereoSine(left, right int, sampleRate int) (s *StereoSine, err error) {
	s = &StereoSine{}
	if s.stream, err = portaudio.OpenDefaultStream(0, 2, float64(sampleRate), 0, s.process); err != nil {
		return
	}

	s.Left.sampleRate = float64(sampleRate)
	s.Left.SetFrequency(left)
	s.Left.SetVolume(100)

	s.Right.sampleRate = float64(sampleRate)
	s.Right.SetFrequency(right)
	s.Right.SetVolume(100)

	return
}

type StereoSine struct {
	Left, Right             channel
	leftVolume, rightVolume int
	stream                  *portaudio.Stream
}

func (s *StereoSine) Play() error {
	return s.stream.Start()
}

func (s *StereoSine) Stop() error {
	return s.stream.Stop()
}

func (s *StereoSine) SetFrequency(freq int) {
	s.Left.SetFrequency(freq)
	s.Right.SetFrequency(freq)
}

func (s *StereoSine) SetVolume(volume int) {
	s.Left.SetVolume(volume)
	s.Right.SetVolume(volume)
}

func (s *StereoSine) process(out [][]float32) {
	for i := range out[0] {
		out[0][i] = float32(s.Left.calc())
		out[1][i] = float32(s.Right.calc())
	}
}

type channel struct {
	sampleRate  float64
	freq, phase float64
	volume      float64
}

func (c *channel) calc() (res float64) {
	res = math.Sin(2*math.Pi*c.phase) * c.volume
	_, c.phase = math.Modf(c.phase + c.freq)
	return
}

func (c *channel) SetFrequency(freq int) {
	c.freq = float64(freq) / c.sampleRate
}

func (c *channel) SetVolume(volume int) {
	c.volume = float64(volume) / 100
	if c.volume > 1 {
		c.volume = 1
	}
	if c.volume < 0 {
		c.volume = 0
	}
}
