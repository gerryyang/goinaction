// Copyright © 2011-12 Qtrac Ltd.
// 
// This program or package and any associated files are licensed under the
// Apache License, Version 2.0 (the "License"); you may not use these files
// except in compliance with the License. You can get a copy of the License
// at: http://www.apache.org/licenses/LICENSE-2.0.
// 
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package shapes

import (
    "fmt"
    "image"
    "image/color"
    "image/draw"
    "image/jpeg"
    "image/png"
    "math"
    "os"
    "path/filepath"
    "runtime"
    "strings"
)


func clamp(minimum, x, maximum int) int {
    switch {
    case x < minimum:
        return minimum
    case x > maximum:
        return maximum
    }
    return x
}

func validFillColor(fill color.Color) color.Color {
    if fill == nil { // We silently treat a nil color as black
        return color.Black
    }
    return fill
}

type Drawer interface {
    Draw(img draw.Image, x, y int) error
}

type Circle struct {
    color.Color
    Radius int
}

func (circle Circle) Draw(img draw.Image, x, y int) error {
    // Algorithm taken from
    // http://en.wikipedia.org/wiki/Midpoint_circle_algorithm
    // No need to check the radius is in bounds because you can only
    // create circles using NewCircle() which guarantees it is within
    // bounds. But the x, y might be outside the image so we check.
    if err := checkBounds(img, x, y); err != nil {
        return err
    }
    fill := validFillColor(circle.Color)
    radius := clamp(1, circle.Radius, 1024)

    x0, y0 := x, y
    f := 1 - radius
    ddF_x, ddF_y := 1, -2*radius
    x, y = 0, radius

    img.Set(x0, y0+radius, fill)
    img.Set(x0, y0-radius, fill)
    img.Set(x0+radius, y0, fill)
    img.Set(x0-radius, y0, fill)

    for x < y {
        if f >= 0 {
            y--
            ddF_y += 2
            f += ddF_y
        }
        x++
        ddF_x += 2
        f += ddF_x
        img.Set(x0+x, y0+y, fill)
        img.Set(x0-x, y0+y, fill)
        img.Set(x0+x, y0-y, fill)
        img.Set(x0-x, y0-y, fill)
        img.Set(x0+y, y0+x, fill)
        img.Set(x0-y, y0+x, fill)
        img.Set(x0+y, y0-x, fill)
        img.Set(x0-y, y0-x, fill)
    }
    return nil
}

func (circle Circle) String() string {
    return fmt.Sprintf("circle(fill=%v, radius=%d)", circle.Color,
        circle.Radius)
}

func checkBounds(img image.Image, x, y int) error {
    if !image.Rect(x, y, x, y).In(img.Bounds()) {
        return fmt.Errorf("%s(): point (%d, %d) is outside the image\n",
            caller(1), x, y)
    }
    return nil
}

func caller(steps int) string {
    name := "?"
    if pc, _, _, ok := runtime.Caller(steps + 1); ok {
        name = filepath.Base(runtime.FuncForPC(pc).Name())
    }
    return name
}

type RegularPolygon struct {
    color.Color
    Radius int
    Sides int
}


func (polygon RegularPolygon) Draw(img draw.Image, x, y int) error {
    // No need to check the radius or sides are in bounds because you can
    // only create polygons using NewRegularPolygon() which guarantees they
    // are within bounds. But the x, y might be outside the image so we
    // check. len(points) == sides + 1
    if err := checkBounds(img, x, y); err != nil {
        return err
    }
    fill := validFillColor(polygon.Color)
    radius := clamp(1, polygon.Radius, 1024)
    sides := clamp(3, polygon.Sides, 60)
    points := getPoints(x, y, sides, float64(radius))
    for i := 0; i < sides; i++ { // Draw lines between the apexes
        drawLine(img, points[i], points[i+1], fill)
    }
    return nil
}

func getPoints(x, y, sides int, radius float64) []image.Point {
    points := make([]image.Point, sides+1)
    // Compute the shape's apexes (thanks to Jasmin Blanchette)
    fullCircle := 2 * math.Pi
    x0, y0 := float64(x), float64(y)
    for i := 0; i < sides; i++ {
        θ := float64(float64(i) * fullCircle / float64(sides))
        x1 := x0 + (radius * math.Sin(θ))
        y1 := y0 + (radius * math.Cos(θ))
        points[i] = image.Pt(int(x1), int(y1))
    }
    points[sides] = points[0] // close the shape
    return points
}

// Based on my Perl Image::Base.pm module's line() method 
func drawLine(img draw.Image, start, end image.Point,
    fill color.Color) {
    x0, x1 := start.X, end.X
    y0, y1 := start.Y, end.Y
    Δx := math.Abs(float64(x1 - x0))
    Δy := math.Abs(float64(y1 - y0))
    if Δx >= Δy { // shallow slope
        if x0 > x1 {
            x0, y0, x1, y1 = x1, y1, x0, y0
        }
        y := y0
        yStep := 1
        if y0 > y1 {
            yStep = -1
        }
        remainder := float64(int(Δx/2)) - Δx
        for x := x0; x <= x1; x++ {
            img.Set(x, y, fill)
            remainder += Δy
            if remainder >= 0.0 {
                remainder -= Δx
                y += yStep
            }
        }
    } else { // steep slope
        if y0 > y1 {
            x0, y0, x1, y1 = x1, y1, x0, y0
        }
        x := x0
        xStep := 1
        if x0 > x1 {
            xStep = -1
        }
        remainder := float64(int(Δy/2)) - Δy
        for y := y0; y <= y1; y++ {
            img.Set(x, y, fill)
            remainder += Δx
            if remainder >= 0.0 {
                remainder -= Δy
                x += xStep
            }
        }
    }
}

func (polygon RegularPolygon) String() string {
    return fmt.Sprintf("polygon(fill=%v, radius=%d, sides=%d)",
        polygon.Color, polygon.Radius, polygon.Sides)
}

type Option struct {
    Fill   color.Color
    Radius int
}

func New(shape string, option Option) (Drawer, error) {
    sidesForShape := map[string]int{"triangle": 3, "square": 4,
        "pentagon": 5, "hexagon": 6, "heptagon": 7, "octagon": 8,
        "enneagon": 9, "nonagon": 9, "decagon": 10}
    if sides, found := sidesForShape[shape]; found {
        return RegularPolygon{option.Fill, option.Radius, sides}, nil
    }
    if shape != "circle" {
        return nil, fmt.Errorf("shapes.New(): invalid shape '%s'", shape)
    }
    return Circle{option.Fill, option.Radius}, nil
}

func FilledImage(width, height int, fill color.Color) draw.Image {
    fill = validFillColor(fill)
    width = clamp(1, width, 4096)
    height = clamp(1, height, 4096)
    img := image.NewRGBA(image.Rect(0, 0, width, height))
    draw.Draw(img, img.Bounds(), &image.Uniform{fill}, image.ZP, draw.Src)
    return img
}

func DrawShapes(img draw.Image, x, y int, drawers ...Drawer) error {
    for _, drawer := range drawers {
        if err := drawer.Draw(img, x, y); err != nil {
            return err
        }
        // Thicker so that it shows up better in screenshots
        if err := drawer.Draw(img, x+1, y); err != nil {
            return err
        }
        if err := drawer.Draw(img, x, y+1); err != nil {
            return err
        }
    }
    return nil
}

func SaveImage(img image.Image, filename string) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()
    switch strings.ToLower(filepath.Ext(filename)) {
    case ".jpg", ".jpeg":
        return jpeg.Encode(file, img, nil)
    case ".png":
        return png.Encode(file, img)
    }
    return fmt.Errorf("shapes.SaveImage(): '%s' has an unrecognized "+
        "suffix", filename)
}
