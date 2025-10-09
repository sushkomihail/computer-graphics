package lab2

import base "go-graphics/lab1"

type PlaneEquation struct {
	A float32
	B float32
	C float32
	D float32
}

func (p PlaneEquation) GetPlaneZ(x, y float32) float32 {
	return -(p.A*x + p.B*y + p.D) / p.C
}

func GetPlaneEquation(p1, p2, p3 base.Vector) (eq PlaneEquation) {
	vec1 := p2
	vec1.Add(p1.Reverse())

	vec2 := p3
	vec2.Add(p1.Reverse())

	normal := base.Cross(vec1, vec2)
	eq.A, eq.B, eq.C = normal.X, normal.Y, normal.Z

	eq.D = -eq.A*p1.X - eq.B*p1.Y - eq.C*p1.Z
	return eq
}
