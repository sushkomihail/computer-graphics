package lab1

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Object struct {
	vertices                []Vector
	transformedVertices     []Vector
	oldTransformationMatrix *[4][4]float32

	position Vector
	rotation Vector
	scale    Vector
}

func NewObject(vertices []Vector) *Object {
	transformedVertices := make([]Vector, len(vertices))
	oldTransformationMatrix := GetIdentityMatrix()

	copy(transformedVertices, vertices)

	return &Object{
		vertices:                vertices,
		transformedVertices:     transformedVertices,
		oldTransformationMatrix: &oldTransformationMatrix,
		scale:                   NewVector(1, 1, 1),
	}
}

func (o Object) DrawObject(screen *ebiten.Image) {
	objectColor := color.RGBA{255, 255, 255, 255}

	for i := 0; i < len(o.transformedVertices)-5; i++ {
		v0 := o.transformedVertices[i]
		v1 := o.transformedVertices[i+1]
		vector.StrokeLine(screen, v0.X, v0.Y, v1.X, v1.Y, float32(strokeWidth), objectColor, false)
	}

	v0 := o.transformedVertices[0]
	v1 := o.transformedVertices[len(o.transformedVertices)-5]
	vector.StrokeLine(screen, v0.X, v0.Y, v1.X, v1.Y, float32(strokeWidth), objectColor, false)
}

func (o *Object) ApplyTransformation(deltaPosition, deltaRotation, scaleFactor Vector,
	projectionMatrix [4][4]float32, axes CoordinatesSystem) {
	copy(o.transformedVertices, o.vertices)

	o.rotation.Add(deltaRotation)
	o.scale.X *= scaleFactor.X
	o.scale.Y *= scaleFactor.Y
	o.scale.Z *= scaleFactor.Z

	translationMatrix := GetTranslationMatrix(deltaPosition)

	xRotationMatrix := GetRotationXMatrix(deltaRotation.X)
	yRotationMatrix := GetRotationYMatrix(deltaRotation.Y)
	zRotationMatrix := GetRotationZMatrix(deltaRotation.Z)

	scaleMatrix := GetScaleMatrix(scaleFactor)

	resultMatrix := MultiplyMatrices(*o.oldTransformationMatrix, scaleMatrix)
	resultMatrix = MultiplyMatrices(resultMatrix, xRotationMatrix)
	resultMatrix = MultiplyMatrices(resultMatrix, yRotationMatrix)
	resultMatrix = MultiplyMatrices(resultMatrix, zRotationMatrix)
	resultMatrix = MultiplyMatrices(resultMatrix, translationMatrix)

	o.position = NewVector(resultMatrix[0][3], resultMatrix[1][3], resultMatrix[2][3])

	*o.oldTransformationMatrix = resultMatrix

	for i, v := range o.transformedVertices {
		v.ApplyTransformationMatrix4x4(resultMatrix)

		w := v.ApplyTransformationMatrix4x4(projectionMatrix)
		v.X /= w
		v.Y /= w
		v.Z /= w

		axes.ProjectVertex(&v)

		o.transformedVertices[i] = v
	}
}
