package paint

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func AddLabel(img *image.RGBA, x, y int, label string, col color.Color) {
	point := fixed.Point26_6{
		X: fixed.I(x),
		Y: fixed.I(y),
	}
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}

func DrawRectangle(img *image.RGBA, x1, y1, x2, y2 int, clr color.RGBA, thickness int, dotted bool) {
	DrawLine(img, x1, y1, x2, y1, clr, thickness, dotted)
	DrawLine(img, x2, y1, x2, y2, clr, thickness, dotted)
	DrawLine(img, x2, y2, x1, y2, clr, thickness, dotted)
	DrawLine(img, x1, y2, x1, y1, clr, thickness, dotted)
	corner := color.RGBA{68, 68, 68, 255}
	DrawFilledCircle(img, corner, thickness/2, Coords{x1, y1})
	DrawFilledCircle(img, corner, thickness/2, Coords{x2, y1})
	DrawFilledCircle(img, corner, thickness/2, Coords{x2, y2})
	DrawFilledCircle(img, corner, thickness/2, Coords{x1, y2})
}

func DrawRotatedRectangle(img *image.RGBA, x1, y1, x2, y2 int, theta float64, clr color.RGBA, thickness int, dotted bool) {
	sqr := RotateRectangle(x1, y1, x2, y2, theta)

	DrawLine(img, sqr[0][0], sqr[0][1], sqr[1][0], sqr[1][1], clr, thickness, dotted)
	DrawLine(img, sqr[1][0], sqr[1][1], sqr[2][0], sqr[2][1], clr, thickness, dotted)
	DrawLine(img, sqr[2][0], sqr[2][1], sqr[3][0], sqr[3][1], clr, thickness, dotted)
	DrawLine(img, sqr[3][0], sqr[3][1], sqr[0][0], sqr[0][1], clr, thickness, dotted)
	fmt.Println(sqr[0], sqr[1], sqr[2], sqr[3])
	corner := color.RGBA{68, 68, 68, 255}
	DrawFilledCircle(img, corner, thickness/2, Coords{sqr[0][0], sqr[0][1]})
	DrawFilledCircle(img, corner, thickness/2, Coords{sqr[1][0], sqr[1][1]})
	DrawFilledCircle(img, corner, thickness/2, Coords{sqr[2][0], sqr[2][1]})
	DrawFilledCircle(img, corner, thickness/2, Coords{sqr[3][0], sqr[3][1]})
}

// Expects x1,y1 to be top left and x2,y2 to be bottom right
func RotateRectangle(x1, y1, x2, y2 int, theta float64) [4][2]int {
	center := Coords{(x1 + x2) / 2, (y1 + y2) / 2}

	abcd := [4][2]int{
		{x1, y1},
		{x2, y1},
		{x2, y2},
		{x1, y2},
	}
	abcd[0][0], abcd[0][1] = rotatePoint(abcd[0][0], abcd[0][1], center.X, center.Y, theta)
	abcd[1][0], abcd[1][1] = rotatePoint(abcd[1][0], abcd[1][1], center.X, center.Y, theta)
	abcd[2][0], abcd[2][1] = rotatePoint(abcd[2][0], abcd[2][1], center.X, center.Y, theta)
	abcd[3][0], abcd[3][1] = rotatePoint(abcd[3][0], abcd[3][1], center.X, center.Y, theta)
	return [4][2]int{abcd[0], abcd[1], abcd[2], abcd[3]}
}

func DrawTriangle(img *image.RGBA, x1, y1, x2, y2, x3, y3 int, clr color.RGBA, thickness int, dotted bool) {
	DrawLine(img, x1, y1, x2, y2, clr, thickness, dotted)
	DrawLine(img, x2, y2, x3, y3, clr, thickness, dotted)
	DrawLine(img, x3, y3, x1, y1, clr, thickness, dotted)
}

func DrawArc(img *image.RGBA, x, y, radius int, startAngle, endAngle float64, clr color.RGBA, thickness int) {
	for startAngle < 2*math.Pi {
		startAngle += 2 * math.Pi
	}
	for startAngle > 2*math.Pi {
		startAngle -= 2 * math.Pi
	}
	for endAngle < 2*math.Pi {
		endAngle += 2 * math.Pi
	}
	for endAngle > 2*math.Pi {
		endAngle -= 2 * math.Pi
	}
	for endAngle < startAngle {
		endAngle += 2 * math.Pi
	}
	bounds := make(map[int]*Coords)
	for i := startAngle; i < endAngle; i += math.Pi / (float64(radius)) { // tbd
		fmt.Println(i)
		ex := .25 * float64(radius) * (math.Cos(i))
		ey := .25 * float64(radius) * (math.Sin(2*math.Pi - i))
		eyUpper := .25 * float64(radius) * (math.Sin(i))
		rx, ry := int(ex)+x, int(ey)+y
		ryUpper := int(eyUpper) + y
		if rx > img.Bounds().Dx() {
			rx = img.Bounds().Dx()
		}
		if ryUpper > img.Bounds().Dy() {
			ryUpper = img.Bounds().Dy()
		}
		bounds[rx] = &Coords{ry, ryUpper}
		img.Set(rx, ry, clr)
		DrawFilledCircle(img, clr, thickness/2, Coords{rx, ry})
	}
}

func rotatePoint(x, y, cx, cy int, theta float64) (int, int) {
	// Translate point to origin
	x -= cx
	y -= cy
	nx := float64(x)*math.Cos(theta) - float64(y)*math.Sin(theta)
	ny := float64(x)*math.Sin(theta) + float64(y)*math.Cos(theta)
	return int(nx) + cx, int(ny) + cy
}

func DrawLine(img *image.RGBA, x1, y1, x2, y2 int, color color.RGBA, thickness int, dotted bool) {
	DrawFilledCircle(img, color, thickness/2, Coords{x1, y1})
	if x2 < x1 {
		x1, x2 = x2, x1
		y1, y2 = y2, y1
	}
	var slope float64
	var draws int
	diff, diffy := .1, 1.
	if dotted {
		diff, diffy = 7., 7.
	}

	if x2 == x1 {
		if y2 < y1 {
			y1, y2 = y2, y1
		}
		for y := float64(y1); y <= float64(y2); y += diffy {
			DrawFilledCircle(img, color, thickness/2, Coords{x1, int(y)})
			draws++
		}
	} else {
		slope = float64(y2-y1) / float64(x2-x1)
		for x := float64(x1); x <= float64(x2); x += diff {
			y := int(slope*float64(x-float64(x1))) + y1
			DrawFilledCircle(img, color, thickness/2, Coords{int(x), y})
			draws++
		}
	}
	fmt.Println("draws:", draws)
}

func DrawFilledCircle(img *image.RGBA, color color.RGBA, radius int, center Coords) {
	bounds := circleBounds(img, radius, center)
	for x, yBounds := range bounds {
		for yValue := yBounds.X; yValue < yBounds.Y; yValue++ {
			img.Set(x, yValue, color)
		}
	}
}

func DrawCircle(img *image.RGBA, clr color.RGBA, radius int, center Coords) {
	bounds := circleBounds(img, radius, center)
	// DrawFilledCircle(img, clr, radius, center)
	// DrawFilledCircle(img, color.RGBA{}, radius-1, center)
	for x, yBounds := range bounds {
		yBoundsNext, ok := bounds[x+1]
		yBoundsPrev, ok2 := bounds[x-1]
		if ok {
			for y := yBounds.Y; y <= yBoundsNext.Y; y++ {
				img.Set(x, y, clr)
			}
			for y := yBounds.X; y >= yBoundsNext.X; y-- {
				img.Set(x, y, clr)
			}
		}
		if ok2 {
			for y := yBounds.X; y >= yBoundsPrev.X; y-- {
				img.Set(x, y, clr)
			}
			for y := yBounds.Y; y <= yBoundsPrev.Y; y++ {
				img.Set(x, y, clr)
			}
		}
		if ok && !ok2 {
			img.Set(x, yBounds.X, color.RGBA{})
		}
		fmt.Println(yBoundsNext, yBoundsPrev, ok, ok2)
	}
}

type Coords struct {
	X, Y int
}

func circleBounds(img *image.RGBA, radius int, center Coords) map[int]*Coords {
	bounds := make(map[int]*Coords)
	for i := 0.; i < math.Pi; i += math.Pi / (float64(radius)) { // tbd
		ex := .25 * float64(radius) * (math.Cos(i))
		ey := .25 * float64(radius) * (math.Sin(2*math.Pi - i))
		eyUpper := .25 * float64(radius) * (math.Sin(i))
		rx, ry := int(ex)+center.X, int(ey)+center.Y
		ryUpper := int(eyUpper) + center.Y
		if rx > img.Bounds().Dx() {
			rx = img.Bounds().Dx()
		}
		if ryUpper > img.Bounds().Dy() {
			ryUpper = img.Bounds().Dy()
		}
		bounds[rx] = &Coords{ry, ryUpper}
	}
	return bounds
}
