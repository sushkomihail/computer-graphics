package lab1

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

type Point3D struct {
	X, Y, Z float64
}

type Edge struct {
	Start, End int
}

type Model struct {
	Vertices []Point3D
	Edges    []Edge
}

type Transform [4][4]float64

type Game struct {
	model    Model
	rotation Point3D
	position Point3D
	scale    Point3D

	lastMouseX, lastMouseY int
	isRotating             bool
}

func createLetterPModel() Model {
	vertices := []Point3D{
		// Основание буквы П
		{-1.0, 2.0, 0.0}, {1.0, 2.0, 0.0},
		{1.0, -2.0, 0.0}, {-1.0, -2.0, 0.0},

		// Верхняя перекладина
		{-1.0, 0.5, 0.0}, {1.0, 0.5, 0.0},

		// Боковые стороны для объема (Z-координата)
		{-1.0, 2.0, 0.5}, {1.0, 2.0, 0.5},
		{1.0, -2.0, 0.5}, {-1.0, -2.0, 0.5},
		{-1.0, 0.5, 0.5}, {1.0, 0.5, 0.5},
	}

	edges := []Edge{
		// Передняя грань
		{0, 1}, {1, 2}, {2, 3}, {3, 0},
		{4, 5}, {0, 4}, {1, 5}, {2, 5}, {3, 4},

		// Задняя грань
		{6, 7}, {7, 8}, {8, 9}, {9, 6},
		{10, 11}, {6, 10}, {7, 11}, {8, 11}, {9, 10},

		// Соединения
		{0, 6}, {1, 7}, {2, 8}, {3, 9}, {4, 10}, {5, 11},
	}

	return Model{Vertices: vertices, Edges: edges}
}

func identityMatrix() Transform {
	return Transform{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
}

func rotationXMatrix(angle float64) Transform {
	cos := math.Cos(angle)
	sin := math.Sin(angle)
	return Transform{
		{1, 0, 0, 0},
		{0, cos, -sin, 0},
		{0, sin, cos, 0},
		{0, 0, 0, 1},
	}
}

func rotationYMatrix(angle float64) Transform {
	cos := math.Cos(angle)
	sin := math.Sin(angle)
	return Transform{
		{cos, 0, sin, 0},
		{0, 1, 0, 0},
		{-sin, 0, cos, 0},
		{0, 0, 0, 1},
	}
}

func rotationZMatrix(angle float64) Transform {
	cos := math.Cos(angle)
	sin := math.Sin(angle)
	return Transform{
		{cos, -sin, 0, 0},
		{sin, cos, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
}

func translationMatrix(x, y, z float64) Transform {
	return Transform{
		{1, 0, 0, x},
		{0, 1, 0, y},
		{0, 0, 1, z},
		{0, 0, 0, 1},
	}
}

func scaleMatrix(x, y, z float64) Transform {
	return Transform{
		{x, 0, 0, 0},
		{0, y, 0, 0},
		{0, 0, z, 0},
		{0, 0, 0, 1},
	}
}

func multiplyMatrices(a, b Transform) Transform {
	var result Transform
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			result[i][j] = 0
			for k := 0; k < 4; k++ {
				result[i][j] += a[i][k] * b[k][j]
			}
		}
	}
	return result
}

func transformPoint(point Point3D, transform Transform) Point3D {
	x := transform[0][0]*point.X + transform[0][1]*point.Y + transform[0][2]*point.Z + transform[0][3]
	y := transform[1][0]*point.X + transform[1][1]*point.Y + transform[1][2]*point.Z + transform[1][3]
	z := transform[2][0]*point.X + transform[2][1]*point.Y + transform[2][2]*point.Z + transform[2][3]

	return Point3D{X: x, Y: y, Z: z}
}

// Получаем матрицу вращения для объекта
func (g *Game) getRotationMatrix() Transform {
	// Правильный порядок для локальных вращений: X * Y * Z
	rotX := rotationXMatrix(g.rotation.X)
	rotY := rotationYMatrix(g.rotation.Y)
	rotZ := rotationZMatrix(g.rotation.Z)

	// Для локальных вращений: сначала вращаем вокруг X, потом Y, потом Z
	rotation := multiplyMatrices(rotY, rotZ)
	rotation = multiplyMatrices(rotX, rotation)
	return rotation
}

// Преобразуем направление из локальных в мировые координаты
func (g *Game) localDirectionToWorld(local Point3D) Point3D {
	rotation := g.getRotationMatrix()

	// Применяем только вращение (без перемещения)
	dir := transformPoint(local, rotation)

	// Нормализуем длину (чтобы движение было с постоянной скоростью)
	length := math.Sqrt(dir.X*dir.X + dir.Y*dir.Y + dir.Z*dir.Z)
	if length > 0 {
		dir.X /= length
		dir.Y /= length
		dir.Z /= length
	}

	return dir
}

// Полная матрица преобразования модели
func (g *Game) getModelMatrix() Transform {
	// Порядок: Масштаб -> Вращение -> Перемещение
	scale := scaleMatrix(g.scale.X, g.scale.Y, g.scale.Z)
	rotation := g.getRotationMatrix()
	translate := translationMatrix(g.position.X, g.position.Y, g.position.Z)

	// S * R * T
	model := multiplyMatrices(scale, rotation)
	model = multiplyMatrices(model, translate)
	return model
}

func projectPoint(point Point3D) (float64, float64) {
	// Простая перспективная проекция
	distance := 5.0
	factor := 200.0 / (distance + point.Z)
	x := point.X*factor + float64(screenWidth)/2
	y := -point.Y*factor + float64(screenHeight)/2
	return x, y
}

func (g *Game) Update() error {
	// Ограничиваем углы вращения чтобы избежать накопления ошибок
	g.rotation.X = math.Mod(g.rotation.X, 2*math.Pi)
	g.rotation.Y = math.Mod(g.rotation.Y, 2*math.Pi)
	g.rotation.Z = math.Mod(g.rotation.Z, 2*math.Pi)

	// Вращение мышью (вокруг глобальных осей для удобства)
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if g.isRotating {
			dx := float64(x-g.lastMouseX) * 0.01
			dy := float64(y-g.lastMouseY) * 0.01

			g.rotation.Y += dx
			g.rotation.X += dy
		}
		g.lastMouseX, g.lastMouseY = x, y
		g.isRotating = true
	} else {
		g.isRotating = false
	}

	// Перемещение в локальных координатах
	moveSpeed := 0.05
	var localMove Point3D

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		localMove.Z -= moveSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		localMove.Z += moveSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		localMove.X -= moveSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		localMove.X += moveSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		localMove.Y += moveSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		localMove.Y -= moveSpeed
	}

	// Применяем перемещение в локальных координатах
	if localMove.X != 0 || localMove.Y != 0 || localMove.Z != 0 {
		worldMove := g.localDirectionToWorld(localMove)
		g.position.X += worldMove.X
		g.position.Y += worldMove.Y
		g.position.Z += worldMove.Z
	}

	// Вращение вокруг локальных осей клавишами
	rotateSpeed := 0.02
	if ebiten.IsKeyPressed(ebiten.KeyJ) {
		g.rotation.X += rotateSpeed // Вращение вокруг локальной X
	}
	if ebiten.IsKeyPressed(ebiten.KeyK) {
		g.rotation.X -= rotateSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyH) {
		g.rotation.Y += rotateSpeed // Вращение вокруг локальной Y
	}
	if ebiten.IsKeyPressed(ebiten.KeyL) {
		g.rotation.Y -= rotateSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyU) {
		g.rotation.Z += rotateSpeed // Вращение вокруг локальной Z
	}
	if ebiten.IsKeyPressed(ebiten.KeyI) {
		g.rotation.Z -= rotateSpeed
	}

	// Масштабирование
	if ebiten.IsKeyPressed(ebiten.KeyZ) {
		g.scale.X *= 1.01
		g.scale.Y *= 1.01
		g.scale.Z *= 1.01
	}
	if ebiten.IsKeyPressed(ebiten.KeyX) {
		g.scale.X *= 0.99
		g.scale.Y *= 0.99
		g.scale.Z *= 0.99
	}

	// Сброс
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.rotation = Point3D{}
		g.position = Point3D{}
		g.scale = Point3D{X: 1, Y: 1, Z: 1}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{15, 15, 30, 255})

	modelMatrix := g.getModelMatrix()

	// Преобразуем и проецируем вершины
	projected := make([]Point3D, len(g.model.Vertices))
	screenPoints := make([]Point3D, len(g.model.Vertices))

	for i, v := range g.model.Vertices {
		transformed := transformPoint(v, modelMatrix)
		projected[i] = transformed
		x, y := projectPoint(transformed)
		screenPoints[i] = Point3D{X: x, Y: y, Z: transformed.Z}
	}

	// Рисуем ребра
	for _, edge := range g.model.Edges {
		p1, p2 := screenPoints[edge.Start], screenPoints[edge.End]

		// Проверяем что точки перед камерой
		if projected[edge.Start].Z > -4 && projected[edge.End].Z > -4 {
			drawLine(screen, p1.X, p1.Y, p2.X, p2.Y, color.RGBA{50, 180, 255, 255})
		}
	}

	// Информация
	drawSimpleText(screen, 10, 20, "Локальные координаты - Управление:", color.White)
	drawSimpleText(screen, 10, 40, "W/S: Вперед/Назад (лок. Z)", color.White)
	drawSimpleText(screen, 10, 60, "A/D: Влево/Вправо (лок. X)", color.White)
	drawSimpleText(screen, 10, 80, "Q/E: Вверх/Вниз (лок. Y)", color.White)
	drawSimpleText(screen, 10, 100, "H/L: Вращение Y (лок.)", color.White)
	drawSimpleText(screen, 10, 120, "J/K: Вращение X (лок.)", color.White)
	drawSimpleText(screen, 10, 140, "U/I: Вращение Z (лок.)", color.White)
	drawSimpleText(screen, 10, 160, "ЛКМ: Вращение мышью (глоб.)", color.White)
	drawSimpleText(screen, 10, 180, "Z/X: Масштаб, R: Сброс", color.White)
}

func drawLine(screen *ebiten.Image, x1, y1, x2, y2 float64, clr color.Color) {
	dx := math.Abs(x2 - x1)
	dy := math.Abs(y2 - y1)

	if dx == 0 && dy == 0 {
		screen.Set(int(x1), int(y1), clr)
		return
	}

	var steps int
	if dx > dy {
		steps = int(dx)
	} else {
		steps = int(dy)
	}

	if steps == 0 {
		return
	}

	xStep := (x2 - x1) / float64(steps)
	yStep := (y2 - y1) / float64(steps)

	x, y := x1, y1
	for i := 0; i <= steps; i++ {
		ix, iy := int(x), int(y)
		if ix >= 0 && ix < screenWidth && iy >= 0 && iy < screenHeight {
			screen.Set(ix, iy, clr)
		}
		x += xStep
		y += yStep
	}
}

func drawSimpleText(screen *ebiten.Image, x, y int, text string, clr color.Color) {
	// Простая отрисовка текста точками
	for i := 0; i < len(text); i++ {
		if text[i] != ' ' {
			screen.Set(x+i*6, y, clr)
			screen.Set(x+i*6+1, y, clr)
			screen.Set(x+i*6, y+1, clr)
			screen.Set(x+i*6+1, y+1, clr)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func Test() {
	game := &Game{
		model: createLetterPModel(),
		scale: Point3D{X: 1, Y: 1, Z: 1},
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("3D Преобразования - Локальные координаты")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
