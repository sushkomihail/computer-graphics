package lab1

type Animation struct {
	Duration    float64
	ElapsedTime float64
	isFinished  bool
}

func NewAnimation(duration float64) *Animation {
	return &Animation{
		Duration:   duration,
		isFinished: true,
	}
}

func (a *Animation) IsFinished() bool {
	return a.isFinished
}

func (a *Animation) Start() {
	a.ElapsedTime = 0
	a.isFinished = false
}

func (a *Animation) Update(deltaTime float64) {
	a.ElapsedTime += deltaTime

	if a.ElapsedTime >= a.Duration {
		a.ElapsedTime = a.Duration
		a.isFinished = true
	}
}

func AnimatePerspectiveChange(from, to [4][4]float32, animation *Animation) [4][4]float32 {
	progress := animation.ElapsedTime / animation.Duration

	for i := range 4 {
		for j := range 4 {
			from[i][j] += (to[i][j] - from[i][j]) * float32(progress)
		}
	}

	return from
}
