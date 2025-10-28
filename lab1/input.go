package lab1

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	translationSpeed = 2
	rotationSpeed    = 1
	scaleChange      = 0.1
	intensityChange  = 10
	specCoefChange   = 0.1
	diffCoefChange   = 0.1
)

func HandleTranslationInput() Vector {
	deltaPosition := Vector{}

	handlePressedKey(ebiten.KeyW, &deltaPosition.Z, -translationSpeed)
	handlePressedKey(ebiten.KeyS, &deltaPosition.Z, translationSpeed)
	handlePressedKey(ebiten.KeyA, &deltaPosition.X, -translationSpeed)
	handlePressedKey(ebiten.KeyD, &deltaPosition.X, translationSpeed)
	handlePressedKey(ebiten.KeyArrowUp, &deltaPosition.Y, -translationSpeed)
	handlePressedKey(ebiten.KeyArrowDown, &deltaPosition.Y, translationSpeed)

	return deltaPosition
}

type Mouse struct {
	lastCursorPosition Vector
	isDragging         bool
}

func HandleRotationInput(mouse *Mouse) Vector {
	deltaRotation := Vector{}

	handlePressedKey(ebiten.KeyX, &deltaRotation.X, rotationSpeed)
	handlePressedKey(ebiten.KeyY, &deltaRotation.Y, rotationSpeed)
	handlePressedKey(ebiten.KeyZ, &deltaRotation.Z, rotationSpeed)

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		curX, curY := ebiten.CursorPosition()

		if mouse.isDragging {
			dx := float64(curX - int(mouse.lastCursorPosition.X))
			dy := float64(curY - int(mouse.lastCursorPosition.Y))
			deltaRotation.Y += float32(dx * rotationSpeed)
			deltaRotation.X += float32(-dy * rotationSpeed)
		}

		mouse.lastCursorPosition = NewVector(float32(curX), float32(curY), 0)
		mouse.isDragging = true
	} else {
		mouse.isDragging = false
	}

	return deltaRotation
}

func HandleScaleInput() Vector {
	scaleFactor := NewVector(1, 1, 1)
	_, y := ebiten.Wheel()

	if y == 0 {
		return scaleFactor
	}

	if y > 0 {
		scaleFactor.Multiply(1 + scaleChange)
	} else {
		scaleFactor.Multiply(1 - scaleChange)
	}

	return scaleFactor
}

func HandleAnimationInput(animation *Animation) {
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		if animation.IsFinished() {
			animation.Start()
		}
	}
}

func HandleLightingInput(isSimple *bool, intensity, specCoef, diffCoef *float64) {
	if inpututil.IsKeyJustPressed(ebiten.KeyG) {
		*isSimple = !(*isSimple)
	}

	if isComboPressed(ebiten.KeyShift, ebiten.KeyI) {
		updateValue(intensity, -intensityChange, 0, 250)
	} else if inpututil.IsKeyJustPressed(ebiten.KeyI) {
		updateValue(intensity, intensityChange, 0, 250)
	}

	if isComboPressed(ebiten.KeyShift, ebiten.KeyS) {
		updateValue(specCoef, -specCoefChange, 0, 1)
	} else if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		updateValue(specCoef, specCoefChange, 0, 1)
	}

	if isComboPressed(ebiten.KeyShift, ebiten.KeyD) {
		updateValue(diffCoef, -diffCoefChange, 0, 1)
	} else if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		updateValue(diffCoef, diffCoefChange, 0, 1)
	}
}

func handlePressedKey(key ebiten.Key, dest *float32, src float32) {
	if ebiten.IsKeyPressed(key) {
		*dest = src
	}

	if inpututil.IsKeyJustReleased(key) {
		*dest = 0
	}
}

func isComboPressed(modifier ebiten.Key, key ebiten.Key) bool {
	return ebiten.IsKeyPressed(modifier) && inpututil.IsKeyJustPressed(key)
}

func updateValue(value *float64, change, min, max float64) {
	if change < 0 {
		*value += math.Max(-(*value - min), change)
	} else {
		*value += math.Min((max - *value), change)
	}
}
