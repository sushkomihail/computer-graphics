package lab3

import base "go-graphics/lab1"

func animateCameraMoving(from, to base.Vector, animation *base.Animation) base.Vector {
	progress := animation.ElapsedTime / animation.Duration
	moveDirection := to
	moveDirection.Add(from.Reverse())
	moveDirection.Multiply(float32(progress))
	from.Add(moveDirection)
	return from
}
