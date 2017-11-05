package server

import (
	"encoding/json"
	"log"
	"net/http"
)

func Init() (err error) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			json.NewEncoder(w).Encode(&state)
		case http.MethodPost:
			var changed State
			defer r.Body.Close()
			json.NewDecoder(r.Body).Decode(&changed)
			log.Println(changed)
			json.NewEncoder(w).Encode(&changed)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	})
	return http.ListenAndServe(":8080", nil)
}
