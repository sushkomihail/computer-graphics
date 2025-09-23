package lab1

import "math"

func GetIdentityMatrix() [4][4]float32 {
	return [4][4]float32{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
}

func GetTranslationMatrix(vec Vector) [4][4]float32 {
	return [4][4]float32{
		{1, 0, 0, vec.X},
		{0, 1, 0, vec.Y},
		{0, 0, 1, vec.Z},
		{0, 0, 0, 1},
	}
}

func GetRotationXMatrix(ang float32) [4][4]float32 {
	sin, cos := getSinCos(ang)

	return [4][4]float32{
		{1, 0, 0, 0},
		{0, cos, -sin, 0},
		{0, sin, cos, 0},
		{0, 0, 0, 1},
	}
}

func GetRotationYMatrix(ang float32) [4][4]float32 {
	sin, cos := getSinCos(ang)

	return [4][4]float32{
		{cos, 0, sin, 0},
		{0, 1, 0, 0},
		{-sin, 0, cos, 0},
		{0, 0, 0, 1},
	}
}

func GetRotationZMatrix(ang float32) [4][4]float32 {
	sin, cos := getSinCos(ang)

	return [4][4]float32{
		{cos, -sin, 0, 0},
		{sin, cos, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
}

func GetScaleMatrix(change Vector) [4][4]float32 {
	return [4][4]float32{
		{change.X, 0, 0, 0},
		{0, change.Y, 0, 0},
		{0, 0, change.Z, 0},
		{0, 0, 0, 1},
	}
}

func GetPerspectiveProjectionMatrix(focusDistance float32) [4][4]float32 {
	return [4][4]float32{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 0, 0},
		{0, 0, -1 / focusDistance, 1},
	}
}

func MultiplyMatrices(a, b [4][4]float32) [4][4]float32 {
	var res [4][4]float32

	for i := range 4 {
		for j := range 4 {
			res[i][j] = a[i][0]*b[0][j] + a[i][1]*b[1][j] + a[i][2]*b[2][j] + a[i][3]*b[3][j]
		}
	}

	return res
}

func getSinCos(ang float32) (float32, float32) {
	angInRad := ang * math.Pi / 180
	sin, cos := math.Sincos(float64(angInRad))
	return float32(sin), float32(cos)
}
