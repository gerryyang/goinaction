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

package oslice

import "strings"

func New(less func(interface{}, interface{}) bool) *Slice {
    return &Slice{less: less}
}

func NewStringSlice() *Slice {
    return &Slice{less: func(a, b interface{}) bool {
        return a.(string) < b.(string)
    }}
}

func NewCaseFoldedSlice() *Slice {
    return &Slice{less: func(a, b interface{}) bool {
        return strings.ToLower(a.(string)) < strings.ToLower(b.(string))
    }}
}

func NewIntSlice() *Slice {
    return &Slice{less: func(a, b interface{}) bool {
        return a.(int) < b.(int)
    }}
}

type Slice struct {
    slice []interface{}
    less  func(interface{}, interface{}) bool
}

func (slice *Slice) Clear() {
    slice.slice = nil
}

func (slice *Slice) Add(x interface{}) {
    if slice.slice == nil {
        slice.slice = []interface{}{x}
    } else if index := bisectLeft(slice.slice, slice.less, x);
        index == len(slice.slice) {
        slice.slice = append(slice.slice, x)
    } else { // See Chapter 4's InsertStringSlice for the logic
        updatedSlice := make([]interface{}, len(slice.slice)+1)
        at := copy(updatedSlice, slice.slice[:index])
        at += copy(updatedSlice[at:], []interface{}{x})
        copy(updatedSlice[at:], slice.slice[index:])
        slice.slice = updatedSlice
    }
}

func (slice *Slice) Remove(x interface{}) bool {
    index := bisectLeft(slice.slice, slice.less, x)
    for ; index < len(slice.slice); index++ {
        if !slice.less(slice.slice[index], x) &&
            !slice.less(x, slice.slice[index]) {
            slice.slice = append(slice.slice[:index],
                slice.slice[index+1:]...)
            return true
        }
    }
    return false
}

func (slice *Slice) Index(x interface{}) int {
    index := bisectLeft(slice.slice, slice.less, x)
    if index >= len(slice.slice) ||
        slice.less(slice.slice[index], x) ||
        slice.less(x, slice.slice[index]) {
        return -1
    }
    return index
}

func (slice *Slice) At(index int) interface{} {
    return slice.slice[index]
}

func (slice *Slice) Len() int {
    return len(slice.slice)
}

// Return's the index position where x belongs in the slice
func bisectLeft(slice []interface{},
    less func(interface{}, interface{}) bool, x interface{}) int {
    left, right := 0, len(slice)
    for left < right {
        middle := int((left + right) / 2)
        if less(slice[middle], x) {
            left = middle + 1
        } else {
            right = middle
        }
    }
    return left
}
