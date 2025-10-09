package lab1

import "fmt"

type Vector struct {
	X float32
	Y float32
	Z float32
}

func NewVector(x, y, z float32) Vector {
	return Vector{
		X: x,
		Y: y,
		Z: z,
	}
}

func (v *Vector) Set(x, y, z float32) {
	v.X = x
	v.Y = y
	v.Z = z
}

func (v *Vector) ApplyTransformationMatrix4x4(mat [4][4]float32) float32 {
	vec := make([]float32, 4)

	for i := range 4 {
		vec[i] = v.X*mat[i][0] + v.Y*mat[i][1] + v.Z*mat[i][2] + mat[i][3]
	}

	v.X, v.Y, v.Z = vec[0], vec[1], vec[2]
	return vec[3]
}

func (v *Vector) Add(vec Vector) {
	v.X += vec.X
	v.Y += vec.Y
	v.Z += vec.Z
}

func (v *Vector) Multiply(val float32) {
	v.X *= val
	v.Y *= val
	v.Z *= val
}

func (v Vector) Reverse() Vector {
	v.Multiply(-1)
	return v
}

func (v Vector) ToString() string {
	return fmt.Sprintf("x:%.1f y:%.1f z:%.1f", v.X, v.Y, v.Z)
}

func Cross(vec1, vec2 Vector) Vector {
	x := vec1.Y*vec2.Z - vec1.Z*vec2.Y
	y := vec1.Z*vec2.X - vec1.X*vec2.Z
	z := vec1.X*vec2.Y - vec1.Y*vec2.X
	return NewVector(x, y, z)
}
