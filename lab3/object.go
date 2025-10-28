package lab3

import (
	base "go-graphics/lab1"
	graphics "go-graphics/lab2"

	"github.com/hajimehoshi/ebiten/v2"
)

type RotationObject struct {
	vertices []*graphics.Vertex
	polygons []graphics.Polygon
}

func NewRotationObject(curve []base.Vector, segmentsNumber int, axesScale float32) *RotationObject {
	curveLength := len(curve)

	for i := range curve {
		curve[i].Multiply(axesScale)
	}

	vertices := make([]*graphics.Vertex, 0, (curveLength-2)*segmentsNumber+2)
	vertices = append(vertices, &graphics.Vertex{Point: curve[0]}, &graphics.Vertex{Point: curve[curveLength-1]})

	deltaAngle := float64(360) / float64(segmentsNumber)

	for i := 0; i < segmentsNumber; i++ {
		ang := deltaAngle * float64(i)
		curveCopy := make([]base.Vector, curveLength-2)
		copy(curveCopy, curve[1:curveLength-1])

		rotationMatrix := base.GetRotationYMatrix(float32(ang))

		for i := range curveCopy {
			curveCopy[i].ApplyTransformationMatrix4x4(rotationMatrix)
		}

		appendPoints(curveCopy, &vertices)
	}

	polygons := triangulate(vertices, curveLength-2, segmentsNumber, axesScale)

	return &RotationObject{
		vertices: vertices,
		polygons: polygons,
	}
}

func appendPoints(src []base.Vector, dst *[]*graphics.Vertex) {
	for _, p := range src {
		*dst = append(*dst, &graphics.Vertex{Point: p})
	}
}

func (o *RotationObject) ApplyTransformation(axes base.CoordinatesSystem, deltaRotation base.Vector) {
	for i := range o.polygons {
		o.polygons[i].ApplyTransformation(axes, deltaRotation)
	}

	o.calculateVerticesNormals()
}

func (o RotationObject) FillScanningLine(screen *ebiten.Image, y int, distanceToLightSource float64,
	frameBuffer *[]float64, zBuffer *[]float32) {

}

func triangulate(vertices []*graphics.Vertex, nonRepeatablePointsInCurve, curvesNum int, axesScale float32) []graphics.Polygon {
	segmentTrianglesNumber := (nonRepeatablePointsInCurve-1)*2 + 2
	trianglesNumber := segmentTrianglesNumber * curvesNum
	triangles := make([]graphics.Polygon, 0, trianglesNumber)

	for i := range curvesNum - 1 {
		segmentCurves := [2][]*graphics.Vertex{
			vertices[i*nonRepeatablePointsInCurve+2 : (i+1)*nonRepeatablePointsInCurve+2],
			vertices[(i+1)*nonRepeatablePointsInCurve+2 : (i+2)*nonRepeatablePointsInCurve+2],
		}
		segmentTriangles := triangulateSegment(vertices[0], vertices[1], segmentCurves, axesScale)
		triangles = append(triangles, segmentTriangles...)
	}

	segmentCurves := [2][]*graphics.Vertex{
		vertices[len(vertices)-nonRepeatablePointsInCurve:],
		vertices[2 : nonRepeatablePointsInCurve+2],
	}
	segmentTriangles := triangulateSegment(vertices[0], vertices[1], segmentCurves, axesScale)
	triangles = append(triangles, segmentTriangles...)

	return triangles
}

func triangulateSegment(first, last *graphics.Vertex, nonRepeatablePoints [2][]*graphics.Vertex, axesScale float32) []graphics.Polygon {
	nonRepeatablePointsInCurve := len(nonRepeatablePoints[0])
	segmentTrianglesNumber := (nonRepeatablePointsInCurve-1)*2 + 2
	triangles := make([]graphics.Polygon, 0, segmentTrianglesNumber)

	for i := 0; i < nonRepeatablePointsInCurve-1; i++ {
		halfQuad1 := *graphics.NewPolygon(
			[]*graphics.Vertex{nonRepeatablePoints[0][i], nonRepeatablePoints[1][i], nonRepeatablePoints[0][i+1]},
			axesScale,
			nil,
		)
		halfQuad2 := *graphics.NewPolygon(
			[]*graphics.Vertex{nonRepeatablePoints[0][i+1], nonRepeatablePoints[1][i], nonRepeatablePoints[1][i+1]},
			axesScale,
			nil,
		)
		triangles = append(triangles, halfQuad1, halfQuad2)
	}

	topTriangle := *graphics.NewPolygon(
		[]*graphics.Vertex{
			first,
			nonRepeatablePoints[1][0],
			nonRepeatablePoints[0][0],
		},
		axesScale,
		nil,
	)
	bottomTriangle := *graphics.NewPolygon(
		[]*graphics.Vertex{
			last,
			nonRepeatablePoints[0][nonRepeatablePointsInCurve-1],
			nonRepeatablePoints[1][nonRepeatablePointsInCurve-1],
		},
		axesScale,
		nil,
	)
	triangles = append(triangles, topTriangle, bottomTriangle)
	return triangles
}

func (o RotationObject) calculateVerticesNormals() {
	for _, v := range o.vertices {
		calculateVertexNormal(v, o.polygons)
	}
}

func calculateVertexNormal(vertex *graphics.Vertex, polygons []graphics.Polygon) {
	normal := base.Vector{}
	neighboursCount := 0

	for _, p := range polygons {
		isNeighbour := false

		for _, v := range p.GetVertices() {
			if v == vertex {
				isNeighbour = true
				break
			}
		}

		if isNeighbour {
			normal.Add(base.NewVector(
				p.GetPlaneEquation().A,
				p.GetPlaneEquation().B,
				p.GetPlaneEquation().C,
			))
			neighboursCount++
		}
	}

	if neighboursCount != 0 {
		normal.Multiply(1 / float32(neighboursCount))
	}

	normal.Normalize()
	vertex.Normal = normal
}
