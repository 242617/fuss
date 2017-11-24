package main

/*func fluidHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, "/fluid")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if fluidEnabled {
		stopCh <- struct{}{}
	}
	fluidEnabled = true

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
	parseBody(w, r, &settings)
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

		freq := cycle(*settings.Frequency.Min, *settings.Frequency.Max)
		for {
			f, _ := freq()
			delta := cycle(*settings.Delta.Min, *settings.Delta.Max)
			for {
				select {
				case <-stopCh:
					return
				default:
				}
				d, changed := delta()
				if changed {
					f, _ = freq()
				}
				ss.Left.SetFrequency(f)
				ss.Right.SetFrequency(f + d)
				time.Sleep(time.Duration(*settings.Delay) * time.Millisecond)
			}
		}

	}()
	ss.Play()

}

func cycle(min, max int) func() (int, bool) {
	if min == max {
		return func() (int, bool) {
			return min, true
		}
	}

	cur := min
	inc := max > min
	frw := true
	fst := true

	return func() (int, bool) {
		var chd bool
		if fst {
			fst = false
			return cur, chd
		}

		if frw {
			cur = incr(cur, 1, inc)
			if cur == max {
				frw = false
			}
		} else {
			cur = incr(cur, -1, inc)
			if cur == min {
				frw = true
				chd = true
			}
		}
		return cur, chd
	}
}

func incr(n, m int, inc bool) int {
	if inc {
		n += m
	} else {
		n -= m
	}
	return n
}*/
