package main

/*func manualHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, "/manual")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if fluidEnabled {
		stopCh <- struct{}{}
	}
	fluidEnabled = false

	var settings struct {
		Left  *int `json:"left"`
		Right *int `json:"right"`
	}
	parseBody(w, r, &settings)
	if settings.Left != nil {
		ss.Left.SetFrequency(*settings.Left)
	}
	if settings.Right != nil {
		ss.Right.SetFrequency(*settings.Right)
	}
	ss.Play()
}*/
