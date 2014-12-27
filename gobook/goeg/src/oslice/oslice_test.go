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

package oslice_test

import (
    "fmt"
    "oslice"
    "reflect"
    "testing"
)

func TestSortedIntList(t *testing.T) {
    slice := oslice.NewIntSlice()
    for _, x := range []int{5, 8, -1, 3, 4, 22} {
        slice.Add(x)
    }
    checkList(slice, []int{-1, 3, 4, 5, 8, 22}, t)
    for _, x := range []int{5, 5, 6} {
        slice.Add(x)
    }
    checkList(slice, []int{-1, 3, 4, 5, 5, 5, 6, 8, 22}, t)
    if slice.Index(4) != 2 {
        t.Fatal("4 missing")
    }
    printSlice(slice)
    if slice.Index(99) != -1 {
        t.Fatal("99 wrongly present")
    }
    if slice.Remove(99) != false {
        t.Fatal("99 wrongly removed")
    }
    if slice.Remove(5) != true {
        t.Fatal("5 not removed")
    }
    checkList(slice, []int{-1, 3, 4, 5, 5, 6, 8, 22}, t)
    if slice.Remove(5) != true {
        t.Fatal("5 not removed")
    }
    checkList(slice, []int{-1, 3, 4, 5, 6, 8, 22}, t)
    if slice.Remove(5) != true {
        t.Fatal("5 not removed")
    }
    checkList(slice, []int{-1, 3, 4, 6, 8, 22}, t)
    if slice.Index(5) != -1 {
        t.Fatalf("5 wrongly present at %d", slice.Index(5))
    }
    printSlice(slice)
    slice.Clear()
    if slice.Len() != 0 {
        t.Fatal("cleared list isn't empty")
    }
    if slice.Remove(9) != false {
        t.Fatal("9 wrongly removed")
    }
    if slice.Index(9) != -1 {
        t.Fatal("9 wrongly found")
    }
}

func printSlice(slice *oslice.Slice) {
    fmt.Print("[")
    sep := ", "
    for i := 0; i < slice.Len(); i++ {
        if i+1 == slice.Len() {
            sep = "]\n"
        }
        fmt.Print(slice.At(i), sep)
    }
}

func checkList(slice *oslice.Slice, ints []int, t *testing.T) {
    if slice.Len() != len(ints) {
        t.Fatalf("slice.Len()=%d != %d", slice.Len(), len(ints))
    }
    for i := 0; i < slice.Len(); i++ {
        if !reflect.DeepEqual(slice.At(i), ints[i]) {
            t.Fatalf("mismtach At(%d) %v vs. %d", i, slice.At(i), ints[i])
        }
    }
}
