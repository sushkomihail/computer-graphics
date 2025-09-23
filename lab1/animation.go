package lab1

type Animation struct {
	duration    float64
	elapsedTime float64
	isFinished  bool
}

func NewAnimation(duration float64) *Animation {
	return &Animation{
		duration:   duration,
		isFinished: true,
	}
}

func (a *Animation) IsFinished() bool {
	return a.isFinished
}

func (a *Animation) Start() {
	a.elapsedTime = 0
	a.isFinished = false
}

func (a *Animation) Update(deltaTime float64) {
	a.elapsedTime += deltaTime

	if a.elapsedTime >= a.duration {
		a.elapsedTime = a.duration
		a.isFinished = true
	}
}

func AnimatePerspectiveChange(from, to [4][4]float32, animation *Animation) [4][4]float32 {
	progress := animation.elapsedTime / animation.duration

	for i := range 4 {
		for j := range 4 {
			from[i][j] += (to[i][j] - from[i][j]) * float32(progress)
		}
	}

	return from
}
