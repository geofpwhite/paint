package paint

import (
	"image"
	"image/color"
	"math"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

type Coords struct {
	X, Y int
}

func AddLabel(img *image.RGBA, pos Coords, label string, col color.Color) {
	point := fixed.Point26_6{
		X: fixed.I(pos.X),
		Y: fixed.I(pos.Y),
	}
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}

func DrawRectangle(img *image.RGBA, p1, p2 Coords, clr color.RGBA, thickness int, dotted bool) {
	DrawLine(img, p1, Coords{p2.X, p1.Y}, clr, thickness, dotted)
	DrawLine(img, Coords{p2.X, p1.Y}, p2, clr, thickness, dotted)
	DrawLine(img, p2, Coords{p1.X, p2.Y}, clr, thickness, dotted)
	DrawLine(img, Coords{p1.X, p2.Y}, p1, clr, thickness, dotted)
	corner := color.RGBA{68, 68, 68, 255}
	DrawFilledCircle(img, corner, thickness/2, p1)
	DrawFilledCircle(img, corner, thickness/2, Coords{p2.X, p1.Y})
	DrawFilledCircle(img, corner, thickness/2, p2)
	DrawFilledCircle(img, corner, thickness/2, Coords{p1.X, p2.Y})
}

func DrawRotatedRectangle(img *image.RGBA, p1, p2 Coords, theta float64, clr color.RGBA, thickness int, dotted bool) {
	sqr := RotateRectangle(p1, p2, theta)

	DrawLine(img, sqr[0], sqr[1], clr, thickness, dotted)
	DrawLine(img, sqr[1], sqr[2], clr, thickness, dotted)
	DrawLine(img, sqr[2], sqr[3], clr, thickness, dotted)
	DrawLine(img, sqr[3], sqr[0], clr, thickness, dotted)
	corner := color.RGBA{68, 68, 68, 255}
	DrawFilledCircle(img, corner, thickness/2, sqr[0])
	DrawFilledCircle(img, corner, thickness/2, sqr[1])
	DrawFilledCircle(img, corner, thickness/2, sqr[2])
	DrawFilledCircle(img, corner, thickness/2, sqr[3])
}

// Expects p1 to be top left and p2 to be bottom right
func RotateRectangle(p1, p2 Coords, theta float64) [4]Coords {
	center := Coords{(p1.X + p2.X) / 2, (p1.Y + p2.Y) / 2}

	abcd := [4]Coords{
		p1,
		{p2.X, p1.Y},
		p2,
		{p1.X, p2.Y},
	}
	abcd[0] = RotatePoint(abcd[0], center, theta)
	abcd[1] = RotatePoint(abcd[1], center, theta)
	abcd[2] = RotatePoint(abcd[2], center, theta)
	abcd[3] = RotatePoint(abcd[3], center, theta)
	return abcd
}

func DrawTriangle(img *image.RGBA, p1, p2, p3 Coords, clr color.RGBA, thickness int, dotted bool) {
	DrawLine(img, p1, p2, clr, thickness, dotted)
	DrawLine(img, p2, p3, clr, thickness, dotted)
	DrawLine(img, p3, p1, clr, thickness, dotted)
}

func DrawArc(img *image.RGBA, center Coords, radius int, startAngle, endAngle float64, clr color.RGBA, thickness int) {
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
		img.Set(rx, ry, clr)
		DrawFilledCircle(img, clr, thickness/2, Coords{rx, ry})
	}
}

func RotatePoint(p, center Coords, theta float64) Coords {
	// Translate point to origin
	x := p.X - center.X
	y := p.Y - center.Y
	nx := float64(x)*math.Cos(theta) - float64(y)*math.Sin(theta)
	ny := float64(x)*math.Sin(theta) + float64(y)*math.Cos(theta)
	return Coords{int(nx) + center.X, int(ny) + center.Y}
}

func DrawLine(img *image.RGBA, p1, p2 Coords, color color.RGBA, thickness int, dotted bool) {
	DrawFilledCircle(img, color, thickness/2, p1)
	if p2.X < p1.X {
		p1, p2 = p2, p1
	}
	var slope float64
	var draws int
	diff, diffy := .1, 1.
	if dotted {
		diff, diffy = 7., 7.
	}

	if p2.X == p1.X {
		if p2.Y < p1.Y {
			p1.Y, p2.Y = p2.Y, p1.Y
		}
		for y := float64(p1.Y); y <= float64(p2.Y); y += diffy {
			DrawFilledCircle(img, color, thickness/2, Coords{p1.X, int(y)})
			draws++
		}
	} else {
		slope = float64(p2.Y-p1.Y) / float64(p2.X-p1.X)
		for x := float64(p1.X); x <= float64(p2.X); x += diff {
			y := int(slope*float64(x-float64(p1.X))) + p1.Y
			DrawFilledCircle(img, color, thickness/2, Coords{int(x), y})
			draws++
		}
	}
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
	}
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
