package lab3

import (
	"math"

	base "go-graphics/lab1"
)

func semicircleFunction(r, x float64) float64 {
	return math.Sqrt(r*r - x*x)
}

func getSemicircle(pointsNumber int, r float64) []base.Vector {
	curve := make([]base.Vector, pointsNumber)
	step := 2 * r / float64(pointsNumber-1)
	y := -r

	for i := range pointsNumber {
		x := semicircleFunction(r, y)
		curve[i] = base.NewVector(float32(x), float32(y), 0)
		y += step
	}

	return curve
}

func generateProfile() []base.Vector {
	var profile []base.Vector
	for i := 0; i <= 20; i++ {
		t := float64(i) / float64(20)
		// Простой профиль для тестирования
		x := 0.2 + 0.3*math.Sin(t*math.Pi)
		y := 2*t - 1
		profile = append(profile, base.NewVector(float32(x), float32(y), 0))
	}
	return profile
}

// func getCone(baseAngle float64) []base.Vector {
// 	a := 1 / math.Tan(baseAngle*math.Pi/180)
// 	return []base.Vector{
// 		base.NewVector(0, -1, 0),
// 		base.NewVector(float32(a), 1, 0),
// 		base.NewVector(0, 1, 0),
// 	}
// }
