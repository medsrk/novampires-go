package common

import "math"

func NormalizeAngle(angle float64) float64 {
	return angle - 2*math.Pi*math.Floor((angle+math.Pi)/(2*math.Pi))
}
