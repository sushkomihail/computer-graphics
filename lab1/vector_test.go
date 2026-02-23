package lab1

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVector_CalculateLength(t *testing.T) {
	testTable := []struct {
		vector   Vector
		expected float64
	}{
		{Vector{}, 0},
		{NewVector(3, 4, 0), 5},
	}

	for _, test := range testTable {
		result := test.vector.CalculateLength()
		assert.Equal(t, test.expected, result,
			fmt.Sprintf("incorrect length for the vector %v, expect %f, got %f", test.vector, test.expected, result))
	}
}

func TestVector_Add(t *testing.T) {
	testTable := []struct {
		vector      Vector
		vectorToAdd Vector
		expected    Vector
	}{
		{Vector{}, NewVector(1, 1, 1), NewVector(1, 1, 1)},
		{NewVector(1, 1, 1), NewVector(-2, 3, -4), NewVector(-1, 4, -3)},
		{NewVector(1, 1, 1), Vector{}, NewVector(1, 1, 1)},
	}

	for _, test := range testTable {
		test.vector.Add(test.vectorToAdd)
		assert.Equal(t, test.expected, test.vector,
			fmt.Sprintf("incorrect vector, expect %v, got %v", test.expected, test.vector))
	}
}

func TestVector_Multiply(t *testing.T) {
	testTable := []struct {
		vector   Vector
		value    float32
		expected Vector
	}{
		{Vector{}, 1, Vector{}},
		{NewVector(1, 2, 3), 5, NewVector(5, 10, 15)},
		{NewVector(-1, 2, -3), -5, NewVector(5, -10, 15)},
	}

	for _, test := range testTable {
		test.vector.Multiply(test.value)
		assert.Equal(t, test.expected, test.vector,
			fmt.Sprintf("incorrect vector, expect %v, got %v", test.expected, test.vector))
	}
}

func TestVector_Reverse(t *testing.T) {
	testTable := []struct {
		vector   Vector
		expected Vector
	}{
		{Vector{}, Vector{}},
		{NewVector(1, 1, 1), NewVector(-1, -1, -1)},
		{NewVector(-1, 2, -3), NewVector(1, -2, 3)},
	}

	for _, test := range testTable {
		result := test.vector.Reverse()
		assert.Equal(t, test.expected, result,
			fmt.Sprintf("incorrect vector, expect %v, got %v", test.expected, result))
	}
}

func TestVector_ComparableWith(t *testing.T) {
	testTable := []struct {
		vector           Vector
		comparableVector Vector
		expected         bool
	}{
		{Vector{}, Vector{}, true},
		{NewVector(1, 1, 0), NewVector(1, 1, 1), false},
	}

	for _, test := range testTable {
		result := test.vector.ComparableWith(test.comparableVector)
		assert.Equal(t, test.expected, result,
			fmt.Sprintf("%v is not comparable to %v", test.vector, test.comparableVector))
	}
}

func TestVector_Dot(t *testing.T) {
	testTable := []struct {
		vector1  Vector
		vector2  Vector
		expected float64
	}{
		{NewVector(1, 0, 0), NewVector(0, 0, 1), 0},
		{NewVector(1, 2, 3), NewVector(4, 5, 6), 32},
	}

	for _, test := range testTable {
		result := Dot(test.vector1, test.vector2)
		assert.Equal(t, test.expected, result,
			fmt.Sprintf("incorrect dot for the vectors: %v, %v, expect %f, got %f",
				test.vector1, test.vector2, test.expected, result))
	}
}

func TestVector_CalculateDistance(t *testing.T) {
	testTable := []struct {
		point1   Vector
		point2   Vector
		expected float64
	}{
		{NewVector(3, 0, 0), NewVector(0, 4, 0), 5},
	}

	for _, test := range testTable {
		result := CalculateDistance(test.point1, test.point2)
		assert.Equal(t, test.expected, result,
			fmt.Sprintf("incorrect distance between points: %v, %v, expect %f, got %f",
				test.point1, test.point2, test.expected, result))
	}
}

func TestVector_Cross(t *testing.T) {
	testTable := []struct {
		vector1  Vector
		vector2  Vector
		expected Vector
	}{
		{NewVector(1, 0, 0), NewVector(0, 0, 1), NewVector(0, -1, 0)},
	}

	for _, test := range testTable {
		result := Cross(test.vector1, test.vector2)
		assert.Equal(t, test.expected, result,
			fmt.Sprintf("incorrect vector, expect %v, got %v", test.expected, result))
	}
}

func TestVector_GetReflected(t *testing.T) {
	testTable := []struct {
		vector   Vector
		normal   Vector
		expected Vector
	}{
		{NewVector(1, 1, 0), NewVector(0, -1, 0), NewVector(-1, 1, 0)},
	}

	for _, test := range testTable {
		result := test.vector.GetReflected(test.normal)
		assert.Equal(t, test.expected, result,
			fmt.Sprintf("incorrect vector, expect %v, got %v", test.expected, result))
	}
}

func TestVector_Normalize(t *testing.T) {
	testTable := []struct {
		vector   Vector
		expected Vector
	}{
		{Vector{}, Vector{}},
		{NewVector(3, 0, 0), NewVector(1, 0, 0)},
	}

	for _, test := range testTable {
		test.vector.Normalize()
		assert.Equal(t, test.expected, test.vector,
			fmt.Sprintf("incorrect vector, expect %v, got %v", test.expected, test.vector))
	}
}
