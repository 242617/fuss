package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gordonklaus/portaudio"

	"github.com/242617/fuss/sine"
	"github.com/242617/utils/parse"
)

const (
	DefaultDeltaMin, DefaultDeltaMax         int = 7, 10
	DefaultFrequencyMin, DefaultFrequencyMax int = 30, 50
	DefaultDelay                             int = 250
)

var (
	err    error
	fluid  bool
	ss     *sine.StereoSine
	stopCh chan struct{}
)

var config struct {
	Address string `json:"address"`
	Static  string `json:"static"`
}

var configPath = flag.String("config", "./fuss.config.json", "config path")

func main() {
	log.SetFlags(log.Lshortfile)

	if err := parse.ParseConfig(*configPath, &config); err != nil {
		log.Fatal(err)
	}

	log.Println("start fuss")
	flag.Parse()

	portaudio.Initialize()
	defer portaudio.Terminate()

	stopCh := make(chan struct{})
	defer close(stopCh)

	if ss, err = sine.NewStereoSine(DefaultFrequencyMin, DefaultFrequencyMin+DefaultDeltaMin, 44100); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, "/check")
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Write([]byte("ok"))
	})
	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, "/stop")
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if fluid {
			stopCh <- struct{}{}
			fluid = false
		}

		ss.Stop()
	})
	http.HandleFunc("/fluid", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, "/fluid")
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if fluid {
			stopCh <- struct{}{}
		}
		fluid = true

		defer r.Body.Close()
		if barr, err := ioutil.ReadAll(r.Body); err != nil {
			log.Println(err)
		} else {
			var settings struct {
				Delta struct {
					Min *int `json:"min"`
					Max *int `json:"max"`
				} `json:"delta"`
				Frequency struct {
					Min *int `json:"min"`
					Max *int `json:"max"`
				} `json:"frequency"`
				Delay *int `json:"delay"`
			}
			if err := json.Unmarshal(barr, &settings); err != nil {
				log.Println(err)
			} else {

				ss.Play()
				go func() {
					if settings.Delta.Min == nil {
						*settings.Delta.Min = DefaultDeltaMin
					}
					if settings.Delta.Max == nil {
						*settings.Delta.Max = DefaultDeltaMax
					}
					if settings.Frequency.Min == nil {
						*settings.Frequency.Min = DefaultFrequencyMin
					}
					if settings.Frequency.Max == nil {
						*settings.Frequency.Max = DefaultFrequencyMax
					}
					if settings.Delay == nil {
						*settings.Delay = DefaultDelay
					}

					for l := range cycleChan(*settings.Frequency.Min, *settings.Frequency.Max) {
						for d := range cycleChan(*settings.Delta.Min, *settings.Delta.Max) {
							select {
							case <-stopCh:
								return
							default:
								ss.SetLeft(l)
								ss.SetRight(l + d)
								time.Sleep(time.Duration(*settings.Delay) * time.Millisecond)
							}
						}
					}

				}()
			}
		}
	})
	http.HandleFunc("/manual", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, "/manual")
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if fluid {
			stopCh <- struct{}{}
			fluid = false
		}

		defer r.Body.Close()
		if barr, err := ioutil.ReadAll(r.Body); err != nil {
			log.Println(err)
		} else {
			var settings struct {
				Left  *int `json:"left"`
				Right *int `json:"right"`
			}
			if err := json.Unmarshal(barr, &settings); err != nil {
				log.Println(err)
			} else {
				if settings.Left != nil {
					ss.SetLeft(*settings.Left)
				}
				if settings.Right != nil {
					ss.SetRight(*settings.Right)
				}
				ss.Play()
			}
		}
	})
	http.Handle("/", http.FileServer(http.Dir(config.Static)))
	log.Fatal(http.ListenAndServe(config.Address, nil))
}

func cycleChan(min, max int) chan int {
	if min > max {
		min, max = max, min
	}
	ch := make(chan int)
	go func() {
		for {
			for i := min; i < max; i++ {
				ch <- i
			}
			for i := max; i > min; i-- {
				ch <- i
			}
		}
	}()
	return ch
}
