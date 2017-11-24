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
