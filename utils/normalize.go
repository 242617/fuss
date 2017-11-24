package utils

func NormalizeVolume(v int) int {
	switch true {
	case v > 100:
		return 100
	case v < 0:
		return 0
	default:
		return v
	}
}

func NormalizeFrequency(f int) int {
	switch true {
	case f > 20000:
		return 100
	case f < 0:
		return 0
	default:
		return f
	}
}
