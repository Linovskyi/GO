package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"sort"

	"github.com/fogleman/gg"
)

func main() {
	imagePath := "C:\\photo.jpg"
	if len(os.Args) > 1 {
		imagePath = os.Args[1]
	}
	regions, err := ExtractRegionsFromImage(imagePath)
	if err != nil {
		panic(err)
	}

	size := len(regions)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			fmt.Printf("%d ", regions[i][j])
		}
		fmt.Println()
	}

	colorCells := make(map[int][][2]int)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			color := regions[i][j]
			colorCells[color] = append(colorCells[color], [2]int{i, j})
		}
	}

	colorOrder := make([]int, 0, len(colorCells))
	for k := range colorCells {
		colorOrder = append(colorOrder, k)
	}
	sort.Slice(colorOrder, func(i, j int) bool {
		return len(colorCells[colorOrder[i]]) < len(colorCells[colorOrder[j]])
	})

	solver := NewColorQueensSolver(size, colorOrder, colorCells)
	solution := solver.Solve()
	if solution == nil {
		fmt.Println("No solution found")
		return
	}

	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if solution[i][j] {
				fmt.Print("Q ")
			} else {
				fmt.Print(". ")
			}
		}
		fmt.Println()
	}

	err = DrawQueens(imagePath, solution)
	if err != nil {
		fmt.Println("Draw error:", err)
	} else {
		fmt.Println("Файл збережено тут:", filepath.Join(".", "result_with_queens.png"))
	}
}

func ExtractRegionsFromImage(imagePath string) ([][]int, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	minCellSize := min(width, height) / 8

	isBlack := func(c color.Color) bool {
		r, g, b, _ := c.RGBA()
		return r>>8 < 80 && g>>8 < 80 && b>>8 < 80
	}

	verticalLines := []int{}
	lastX := -minCellSize
	for x := 0; x < width; x++ {
		blackPixels := 0
		for y := 0; y < height; y++ {
			if isBlack(img.At(x, y)) {
				blackPixels++
			}
		}
		if blackPixels > height/2 {
			if x-lastX > minCellSize/2 {
				verticalLines = append(verticalLines, x)
				lastX = x
			}
		}
	}

	horizontalLines := []int{}
	lastY := -minCellSize
	for y := 0; y < height; y++ {
		blackPixels := 0
		for x := 0; x < width; x++ {
			if isBlack(img.At(x, y)) {
				blackPixels++
			}
		}
		if blackPixels > width/2 {
			if y-lastY > minCellSize/2 {
				horizontalLines = append(horizontalLines, y)
				lastY = y
			}
		}
	}

	fmt.Printf("verticalLines.Count = %d\n", len(verticalLines))
	fmt.Printf("horizontalLines.Count = %d\n", len(horizontalLines))

	rows := len(horizontalLines) - 1
	cols := len(verticalLines) - 1
	regions := make([][]int, rows)
	for i := range regions {
		regions[i] = make([]int, cols)
	}

	colorMap := make(map[[3]uint8]int)
	nextRegion := 0

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			x1, x2 := verticalLines[j], verticalLines[j+1]
			y1, y2 := horizontalLines[i], horizontalLines[i+1]
			cx, cy := (x1+x2)/2, (y1+y2)/2
			r, g, b, _ := img.At(cx, cy).RGBA()
			key := [3]uint8{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8)}
			region, ok := colorMap[key]
			if !ok {
				region = nextRegion
				colorMap[key] = region
				nextRegion++
			}
			regions[i][j] = region
		}
	}
	return regions, nil
}

type ColorQueensSolver struct {
	size       int
	colorOrder []int
	colorCells map[int][][2]int
	board      [][]bool
	rows       []bool
	cols       []bool
}

func NewColorQueensSolver(size int, colorOrder []int, colorCells map[int][][2]int) *ColorQueensSolver {
	board := make([][]bool, size)
	for i := range board {
		board[i] = make([]bool, size)
	}
	return &ColorQueensSolver{
		size:       size,
		colorOrder: colorOrder,
		colorCells: colorCells,
		board:      board,
		rows:       make([]bool, size),
		cols:       make([]bool, size),
	}
}

func (s *ColorQueensSolver) Solve() [][]bool {
	if s.placeQueen(0) {
		return s.board
	}
	return nil
}

func (s *ColorQueensSolver) placeQueen(colorIdx int) bool {
	if colorIdx == len(s.colorOrder) {
		return true
	}
	color := s.colorOrder[colorIdx]
	for _, cell := range s.colorCells[color] {
		i, j := cell[0], cell[1]
		if s.rows[i] || s.cols[j] || s.hasDiagonalNeighbor(i, j) {
			continue
		}
		s.board[i][j] = true
		s.rows[i] = true
		s.cols[j] = true
		if s.placeQueen(colorIdx + 1) {
			return true
		}
		s.board[i][j] = false
		s.rows[i] = false
		s.cols[j] = false
	}
	return false
}

func (s *ColorQueensSolver) hasDiagonalNeighbor(x, y int) bool {
	for dx := -1; dx <= 1; dx += 2 {
		for dy := -1; dy <= 1; dy += 2 {
			nx, ny := x+dx, y+dy
			if nx >= 0 && nx < s.size && ny >= 0 && ny < s.size {
				if s.board[nx][ny] {
					return true
				}
			}
		}
	}
	return false
}

func DrawQueens(imagePath string, solution [][]bool) error {
	file, err := os.Open(imagePath)
	if err != nil {
		return err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	size := len(solution)
	cellWidth := float64(width) / float64(size)
	cellHeight := float64(height) / float64(size)

	dc := gg.NewContext(width, height)
	dc.DrawImage(img, 0, 0)

	fontSize := minf(cellWidth, cellHeight) * 0.7
	// Спробуйте завантажити шрифт з емодзі, якщо він є у вашій системі
	if err := dc.LoadFontFace("Segoe UI Emoji.ttf", fontSize); err != nil {
		dc.LoadFontFace("arial.ttf", fontSize)
	}

	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if solution[i][j] {
				x := float64(j)*cellWidth + cellWidth/2
				y := float64(i)*cellHeight + cellHeight/2
				emoji := "👑"
				w, h := dc.MeasureString(emoji)
				dc.SetColor(color.Black)
				dc.DrawString(emoji, x-w/2, y+h/2)
			}
		}
	}
	return dc.SavePNG("result_with_queens.png")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func minf(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
