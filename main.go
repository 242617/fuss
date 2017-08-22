package main

import (
	"log"
	"time"

	"github.com/gordonklaus/portaudio"

	"github.com/242617/fuss/sine"
)

const (
	delta      int           = 7
	min, max   int           = 30, 70
	sampleRate int           = 44100
	speed      time.Duration = 80
)

var (
	s   *sine.StereoSine
	err error
)

func main() {
	portaudio.Initialize()
	defer portaudio.Terminate()

	if s, err = sine.NewStereoSine(min, min, sampleRate); err != nil {
		log.Fatal(err)
	}
	s.Play()
	defer s.Stop()

	for {
		for d := 7; d < 20; d++ {
			cycle(s, d)
		}
		for d := 20; d > 7; d-- {
			cycle(s, d)
		}
	}
}

func cycle(s *sine.StereoSine, d int) {
	for i := min; i < max; i++ {
		time.Sleep(speed * time.Millisecond)
		s.SetLeft(i)
		s.SetRight(i + d)
	}
	for i := max; i > min; i-- {
		time.Sleep(speed * time.Millisecond)
		s.SetLeft(i)
		s.SetRight(i + d)
	}
}
