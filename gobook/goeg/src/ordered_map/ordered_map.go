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
    "math"
    "qtrac.eu/omap"
    "strings"
)

type Point struct{ X, Y int }

func (point Point) String() string {
    return fmt.Sprintf("(%d, %d)", point.X, point.Y)
}

func main() {
    words := []string{"Puttering", "About", "in", "a", "Small", "Land"}
    wordForWord := omap.NewCaseFoldedKeyed()
    fmt.Println(wordForWord.Len(), "words")
    for _, word := range words {
        wordForWord.Insert(word, strings.ToUpper(word))
    }
    wordForWord.Do(func(key, value interface{}) {
        fmt.Printf("%v→%v\n", key, value)
    })
    fmt.Println("length before deleting:", wordForWord.Len())
    _, containsSmall := wordForWord.Find("small")
    fmt.Println("contains small:", containsSmall)
    for _, key := range []string{"big", "medium", "small"} {
        fmt.Printf("%t ", wordForWord.Delete(key))
    }
    _, containsSmall = wordForWord.Find("small")
    fmt.Println("\nlength after deleting: ", wordForWord.Len())
    fmt.Println("contains small:", containsSmall)

    showMap(wordForWord, words, "words", 20)
    searchMap(wordForWord, "small", "big")

    fmt.Println()

    distanceForPoint := omap.New(func(a, b interface{}) bool {
        α, β := a.(*Point), b.(*Point)
        if α.X != β.X {
            return α.X < β.X
        }
        return α.Y < β.Y
    })
    fmt.Println(distanceForPoint.Len(), "points")
    points := []*Point{{3, 1}, {1, 2}, {2, 3}, {1, 3}, {3, 2}, {2, 1}, {2, 2}}
    for _, point := range points {
        distance := math.Hypot(float64(point.X), float64(point.Y))
        distanceForPoint.Insert(point, distance)
    }
    distanceForPoint.Do(func(key, value interface{}) {
        fmt.Printf("%v → %.2v\n", key, value)
    })
    // No &distanceForPoint because it is already a pointer
    showMap(distanceForPoint, points, "points", 5)
    searchMap(distanceForPoint, &Point{1, 1}, &Point{3, 2})
}

func showMap(omap *omap.Map, data interface{}, name string,
    width int) {
    fmt.Println("original:   ", data)
    fmt.Print("omap keys:   [")
    gap := ""
    omap.Do(func(key, _ interface{}) {
        fmt.Print(gap, key)
        gap = " "
    })
    fmt.Println("]")
    fmt.Print("omap values: [")
    gap = ""
    omap.Do(func(_, value interface{}) {
        fmt.Printf("%s%.*v", gap, width, value)
        gap = " "
    })
    fmt.Println("]")
    fmt.Println(omap.Len(), name)
}

func searchMap(omap *omap.Map, keys ...interface{}) {
    for _, key := range keys {
        if value, found := omap.Find(key); found {
            fmt.Printf("\"%v\" is in the omap with value %v\n", key, value)
        } else {
            fmt.Printf("\"%v\" isn't in the omap\n", key)
        }
    }
}
