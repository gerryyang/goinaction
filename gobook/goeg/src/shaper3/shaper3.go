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

package main

import (
    "fmt"
    "image"
    "image/color"
    "log"
    "os"
    "shaper3/shapes"
)

func main() {
    log.SetFlags(0)
    const width, height = 400, 200
    img := shapes.FilledImage(width, height,
        color.RGBA{0xFF, 0xFF, 0xFF, 0xFF})
    x, y := width/4, height/2

    red := color.RGBA{0xFF, 0, 0, 0xFF}
    blue := color.RGBA{0, 0, 0xFF, 0xFF}
    // Purely for testing New() vs. New*()
    if len(os.Args) == 1 {
        fmt.Println("Using NewCircle() & NewRegularPolygon()")
        circle := shapes.Circle{blue, 90}
        circle.Color = red
        octagon := shapes.RegularPolygon{red, 75, 8}
        octagon.Color = blue
        polygon := shapes.RegularPolygon{image.Black, 65, 4}
        if err := shapes.DrawShapes(img, x, y, circle, octagon, polygon);
            err != nil {
            fmt.Println(err)
        }
        sanityCheck("circle", circle)
        sanityCheck("octagon", octagon)
        sanityCheck("polygon", polygon)
    } else {
        fmt.Println("Using New()")
        if _, err := shapes.New("Misshapen", shapes.Option{blue, 5});
            err == nil {
            fmt.Println("unexpectedly got a non-nil invalid shape!")
        }
        shape, _ := shapes.New("circle", shapes.Option{blue, 5})
        circle := shape.(shapes.Circle)
        circle.Color = red
        circle.Radius = 90
        shape, _ = shapes.New("octagon", shapes.Option{red, 10})
        octagon := shape.(shapes.RegularPolygon)
        octagon.Color = blue
        octagon.Radius = 75
        polygon, _ := shapes.New("square", shapes.Option{Radius: 65})
        if err := shapes.DrawShapes(img, x, y, circle, octagon, polygon);
            err != nil {
            fmt.Println(err)
        }
        sanityCheck("circle", circle)
        sanityCheck("octagon", octagon)
        sanityCheck("polygon", polygon)
    }
    polygon := shapes.RegularPolygon{color.RGBA{0, 0x7F, 0, 0xFF}, 65, 4}
    sanityCheck("polygon", polygon)
    y = 30
    for i, radius := range []int{60, 55, 50, 45, 40} {
        polygon.Radius = radius
        polygon.Sides = i + 5
        x += radius
        y += height / 8
        if err := shapes.DrawShapes(img, x, y, polygon); err != nil {
            fmt.Println(err)
        }
    }

    filename := "shapes.png"
    if err := shapes.SaveImage(img, filename); err != nil {
        log.Println(err)
    } else {
        fmt.Println("Saved", filename)
    }
    fmt.Println("OK")

    img = shapes.FilledImage(width, height, image.White)
    x, y = width/3, height/4
}

func sanityCheck(name string, drawer shapes.Drawer) {
    fmt.Print("name=", name, " ")
    var fill color.Color
    radius, sides := -1, -1
    if circle, ok := drawer.(shapes.Circle); ok {
        fill = circle.Color
        radius = circle.Radius
    }
    if polygon, ok := drawer.(shapes.RegularPolygon); ok {
        fill = polygon.Color
        radius = polygon.Radius
        sides = polygon.Sides
    }
    if fill != nil {
        fmt.Print("fill=", fill, " ")
    }
    if radius != -1 {
        fmt.Print("radius=", radius, " ")
    }
    if sides != -1 {
        fmt.Print("sides=", sides, " ")
    }
    fmt.Println()
}
