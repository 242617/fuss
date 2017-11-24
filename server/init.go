package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/242617/torture/config"
	"github.com/242617/torture/sine"
	"github.com/242617/torture/utils"
)

var start = time.Now()

var ss *sine.StereoSine

type request struct {
	Enabled *bool `json:"enabled"`
	Volume  *int  `json:"volume"`
	Left    *int  `json:"left"`
	Right   *int  `json:"right"`
}

const (
	DefaultDeltaMin, DefaultDeltaMax         int = 7, 10
	DefaultFrequencyMin, DefaultFrequencyMax int = 30, 50
	DefaultDelay                             int = 250
)

func Init() (err error) {

	if ss, err = sine.NewStereoSine(DefaultFrequencyMin, DefaultFrequencyMin+DefaultDeltaMin, 44100); err != nil {
		return
	}

	state = NewState()

	static := http.FileServer(http.Dir(config.Config.Static))

	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		resp := struct {
			Uptime string `json:"uptime"`
		}{time.Since(start).String()}
		json.NewEncoder(w).Encode(&resp)
	})
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		addCORS(&w)
		switch r.Method {
		case http.MethodOptions:
		case http.MethodGet:
			json.NewEncoder(w).Encode(&state)
		case http.MethodPut:
			var changes request
			defer r.Body.Close()
			if err = json.NewDecoder(r.Body).Decode(&changes); err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if err = process(changes); err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(&state)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		addCORS(&w)
		static.ServeHTTP(w, r)
	})
	log.Println("server", fmt.Sprintf(`starting at "%s"`, config.Config.ServerAddress))
	return http.ListenAndServe(config.Config.ServerAddress, nil)
}

func addCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", config.Config.Origin)
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, PUT, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")
}

func process(changes request) (err error) {

	if changes.Enabled != nil {
		log.Println("*changes.Enabled", *changes.Enabled)
		if !state.Enabled && *changes.Enabled {
			ss.Play()
			state.Enabled = true
		}
		if state.Enabled && !*changes.Enabled {
			ss.Stop()
			state.Enabled = false
		}
	}

	if changes.Volume != nil {
		log.Println("*changes.Volume", *changes.Volume)
		state.Volume = utils.NormalizeVolume(*changes.Volume)
		ss.SetVolume(state.Volume)
	}

	if changes.Left != nil {
		log.Println("*changes.Left", *changes.Left)
		state.Left = utils.NormalizeFrequency(*changes.Left)
		ss.Left.SetFrequency(*changes.Left)
	}

	if changes.Right != nil {
		log.Println("*changes.Right", *changes.Right)
		state.Right = utils.NormalizeFrequency(*changes.Right)
		ss.Right.SetFrequency(*changes.Right)
	}

	return
}
