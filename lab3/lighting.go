package lab3

import (
	base "go-graphics/lab1"
	graphics "go-graphics/lab2"
	"math"
)

const (
	bgLightIntensity = 0.3
	bgLightCoef      = 1
	k                = 1
	phongCoef        = 16
)

var (
	lightIntensity = 150.0
	diffReflCoef   = 0.6
	specReflCoef   = 0.0
)

var lightSourcePosition = base.NewVector(windowWidth/2, 0, 150)
var observerPosition = base.NewVector(windowHeight/2, 0, 150)

func calculateIntersectionIntensity(intersection graphics.Intersection, intersectionPoint base.Vector,
	distanceToLightSource float64) float64 {
	lightDirectionA := lightSourcePosition
	lightDirectionA.Add(intersection.GetEdge().A.Point.Reverse())
	lightDirectionB := lightSourcePosition
	lightDirectionB.Add(intersection.GetEdge().B.Point.Reverse())

	observeDirectionA := observerPosition
	observeDirectionA.Add(intersection.GetEdge().A.Point.Reverse())
	observeDirectionB := observerPosition
	observeDirectionB.Add(intersection.GetEdge().B.Point.Reverse())

	intensityA := calculateIntensity(distanceToLightSource, lightDirectionA, intersection.GetEdge().A.Normal, observeDirectionA)
	intensityB := calculateIntensity(distanceToLightSource, lightDirectionB, intersection.GetEdge().B.Normal, observeDirectionB)

	edgeLength := base.CalculateDistance(intersection.GetEdge().A.Point, intersection.GetEdge().B.Point)
	distanceToIntersection := base.CalculateDistance(intersection.GetEdge().A.Point, intersectionPoint)

	var k float64

	if edgeLength != 0 {
		k = distanceToIntersection / edgeLength
	}

	intensity := intensityA + (intensityB-intensityA)*k

	return intensity
}

func calculateIntensity(distance float64, lightDirection, normal, observeDirection base.Vector) float64 {
	reflectedRay := lightDirection.GetReflected(normal)

	lightDirection.Normalize()
	normal.Normalize()
	reflectedRay.Normalize()
	observeDirection.Normalize()

	diffuse := diffReflCoef * base.Dot(normal, lightDirection)
	specular := specReflCoef * math.Pow(base.Dot(reflectedRay, observeDirection), phongCoef)

	intensity := bgLightIntensity*bgLightCoef + lightIntensity/(k+distance)*(diffuse+specular)

	if intensity < 0 {
		intensity = 0
	} else if intensity > 1 {
		intensity = 1
	}

	return intensity
}
