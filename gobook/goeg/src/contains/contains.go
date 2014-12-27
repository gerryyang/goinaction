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
    "reflect"
)

type Slicer interface {
    EqualTo(i int, x interface{}) bool
    Len() int
}

type IntSlice []int

func (slice IntSlice) EqualTo(i int, x interface{}) bool {
    return slice[i] == x.(int)
}
func (slice IntSlice) Len() int { return len(slice) }

func IntIndexSlicer(ints []int, x int) int {
    return IndexSlicer(IntSlice(ints), x)
}

type FloatSlice []float64

func (slice FloatSlice) EqualTo(i int, x interface{}) bool {
    return slice[i] == x.(float64)
}
func (slice FloatSlice) Len() int { return len(slice) }

func FloatIndexSlicer(floats []float64, x float64) int {
    return IndexSlicer(FloatSlice(floats), x)
}

type StringSlice []string

func (slice StringSlice) EqualTo(i int, x interface{}) bool {
    return slice[i] == x.(string)
}
func (slice StringSlice) Len() int { return len(slice) }

func StringIndexSlicer(strs []string, x string) int {
    return IndexSlicer(StringSlice(strs), x)
}

// Returns the index position of x in slice or array xs providing xs's
// items are of the same time as x (integers or strings); returns -1 if x
// isn't in xs. Uses a slow linear search suitable for small amounts of
// unsorted data.
func IndexSlicer(slice Slicer, x interface{}) int {
    for i := 0; i < slice.Len(); i++ {
        if slice.EqualTo(i, x) {
            return i
        }
    }
    return -1
}

// Returns true if x is in slice or array xs providing xs's items are of
// the same time as x (integers or strings). Uses the Index() function
// which does a slow linear search suitable for small amounts of unsorted
// data.
func InSlice(xs interface{}, x interface{}) bool {
    return Index(xs, x) > -1
}

// Returns the index position of x in slice or array xs providing xs's
// items are of the same time as x (integers or strings); returns -1 if x
// isn't in xs. Uses a slow linear search suitable for small amounts of
// unsorted data.
func Index(xs interface{}, x interface{}) int {
    switch slice := xs.(type) {
    case []int:
        for i, y := range slice {
            if y == x.(int) {
                return i
            }
        }
    case []string:
        for i, y := range slice {
            if y == x.(string) {
                return i
            }
        }
    }
    return -1
}

func InSliceReflect(xs interface{}, x interface{}) bool {
    return IndexReflect(xs, x) > -1
}

func IndexReflectX(xs interface{}, x interface{}) int { // Long-winded way
    if slice := reflect.ValueOf(xs); slice.Kind() == reflect.Slice {
        for i := 0; i < slice.Len(); i++ {
            switch y := slice.Index(i).Interface().(type) {
            case int:
                if y == x.(int) {
                    return i
                }
            case string:
                if y == x.(string) {
                    return i
                }
            }
        }
    }
    return -1
}

func IndexReflect(xs interface{}, x interface{}) int {
    if slice := reflect.ValueOf(xs); slice.Kind() == reflect.Slice {
        for i := 0; i < slice.Len(); i++ {
            if reflect.DeepEqual(x, slice.Index(i)) {
                return i
            }
        }
    }
    return -1
}

func IntSliceIndex(xs []int, x int) int {
    for i, y := range xs {
        if x == y {
            return i
        }
    }
    return -1
}

func StringSliceIndex(xs []string, s string) int {
    for i, x := range xs {
        if x == s {
            return i
        }
    }
    return -1
}

func SliceIndex(limit int, predicate func(i int) bool) int {
    for i := 0; i < limit; i++ {
        if predicate(i) {
            return i
        }
    }
    return -1
}

func main() {
    xs := []int{2, 4, 6, 8}
    fmt.Println("5 @", Index(xs, 5), "  6 @", Index(xs, 6))
    ys := []string{"C", "B", "K", "A"}
    fmt.Println("Z @", Index(ys, "Z"), "  A @", Index(ys, "A"))

    fmt.Println("5 @", IndexReflectX(xs, 5), "  6 @", IndexReflectX(xs, 6))
    fmt.Println("Z @", IndexReflectX(ys, "Z"), "  A @",
        IndexReflectX(ys, "A"))
    fmt.Println("5 @", IndexReflect(xs, 5), "  6 @", IndexReflect(xs, 6))
    fmt.Println("Z @", IndexReflect(ys, "Z"), "  A @",
        IndexReflect(ys, "A"))

    fmt.Println("5 @", IntIndexSlicer(xs, 5),
        "  6 @", IntIndexSlicer(xs, 6))
    fmt.Println("Z @", StringIndexSlicer(ys, "Z"),
        "  A @", StringIndexSlicer(ys, "A"))

    sliceIndex()
}

func sliceIndex() {
    xs := []int{2, 4, 6, 8}
    ys := []string{"C", "B", "K", "A"}
    fmt.Println(
        SliceIndex(len(xs), func(i int) bool { return xs[i] == 5 }),
        SliceIndex(len(xs), func(i int) bool { return xs[i] == 6 }),
        SliceIndex(len(ys), func(i int) bool { return ys[i] == "Z" }),
        SliceIndex(len(ys), func(i int) bool { return ys[i] == "A" }))
}
