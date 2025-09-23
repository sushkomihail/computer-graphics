package lab1

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	translationSpeed = 2
	rotationSpeed    = 1
	scaleChange      = 0.1
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

func HandleProjectionSwitchInput(animation *Animation) {
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		if animation.IsFinished() {
			animation.Start()
		}
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
