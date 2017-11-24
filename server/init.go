package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/242617/torture/config"
	"github.com/242617/torture/utils"
)

// var ss sine.StereoSine

var start = time.Now()

type request struct {
	Enabled *bool `json:"enabled,omitempty"`
	Volume  *int  `json:"volume,omitempty"`
}

func Init() (err error) {

	/*if ss, err = sine.NewStereoSine(DefaultFrequencyMin, DefaultFrequencyMin+DefaultDeltaMin, 44100); err != nil {
		log.Fatal(err)
	}*/

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

	log.Println(state, changes)

	if changes.Enabled != nil {
		if !state.Enabled && *changes.Enabled {
			log.Println("start")
			state.Enabled = true
		}
		if state.Enabled && !*changes.Enabled {
			log.Println("stop")
			state.Enabled = false
		}
	}

	if changes.Volume != nil {
		state.Volume = utils.NormalizeVolume(*changes.Volume)
	}

	/*stopCh := make(chan struct{})
	defer close(stopCh)

	ss.Play()
	ss.Left.SetFrequency(80)
	ss.Right.SetFrequency(90)
	ss.SetVolume(100)

	time.Sleep(10 * time.Second)*/

	return

}
