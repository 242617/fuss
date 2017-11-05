package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gordonklaus/portaudio"

	"github.com/242617/torture/config"
	"github.com/242617/torture/server"
	"github.com/242617/torture/sine"
)

const (
	DefaultDeltaMin, DefaultDeltaMax         int = 7, 10
	DefaultFrequencyMin, DefaultFrequencyMax int = 30, 50
	DefaultDelay                             int = 250
)

var (
	err          error
	fluidEnabled bool
	ss           *sine.StereoSine
	stopCh       chan struct{}
)

func main() {
	log.SetFlags(log.Lshortfile)

	flag.StringVar(&config.Path, "config", "torture.yaml", "Application config path")
	flag.Parse()

	if err = config.Init(); err != nil {
		log.Fatal(err)
	}

	log.Fatal(server.Init())

	return

	portaudio.Initialize()
	defer portaudio.Terminate()

	stopCh := make(chan struct{})
	defer close(stopCh)

	if ss, err = sine.NewStereoSine(DefaultFrequencyMin, DefaultFrequencyMin+DefaultDeltaMin, 44100); err != nil {
		log.Fatal(err)
	}

	ss.Play()
	ss.Left.SetFrequency(80)
	ss.Right.SetFrequency(90)
	ss.SetVolume(100)

	time.Sleep(10 * time.Second)

	return

	http.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Method, "/start")
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if fluidEnabled {
			stopCh <- struct{}{}
			fluidEnabled = false
		}

		ss.Play()
	})
	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Method, "/stop")
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if fluidEnabled {
			stopCh <- struct{}{}
			fluidEnabled = false
		}

		ss.Stop()
	})
	http.HandleFunc("/volume", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Method, "/volume")
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var settings struct {
			Left  *int `json:"left"`
			Right *int `json:"right"`
		}
		parseBody(w, r, &settings)
		if settings.Left != nil {
			ss.Left.SetVolume(*settings.Left)
		}
		if settings.Right != nil {
			ss.Right.SetVolume(*settings.Right)
		}
	})
	http.HandleFunc("/fluid", fluidHandler)
	http.HandleFunc("/manual", manualHandler)
	http.HandleFunc("/sms", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Method, "/sms")
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if fluidEnabled {
			stopCh <- struct{}{}
		}
		fluidEnabled = false

		go func() {
			ss.Play()
			ss.Left.SetFrequency(80)
			ss.Right.SetFrequency(90)
			for i := 0; i < 3; i++ {
				ss.SetVolume(100)
				time.Sleep(200 * time.Millisecond)
				ss.SetVolume(0)
				time.Sleep(100 * time.Millisecond)
			}
			ss.Stop()
		}()
	})
	http.HandleFunc("/call", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Method, "/call")
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if fluidEnabled {
			stopCh <- struct{}{}
		}
		fluidEnabled = false

		go func() {
			ss.Play()
			ss.Left.SetFrequency(80)
			ss.Right.SetFrequency(90)
			for i := 0; i < 7; i++ {
				ss.SetVolume(100)
				time.Sleep(800 * time.Millisecond)
				ss.SetVolume(0)
				time.Sleep(1200 * time.Millisecond)
			}
			ss.Stop()
		}()
	})
	http.Handle("/", http.FileServer(http.Dir(config.Config.Static)))
	log.Fatal(http.ListenAndServe(config.Config.ServerAddress, nil))
}

func parseBody(w http.ResponseWriter, r *http.Request, body interface{}) {
	defer r.Body.Close()
	if barr, err := ioutil.ReadAll(r.Body); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	} else {
		if err := json.Unmarshal(barr, &body); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}
