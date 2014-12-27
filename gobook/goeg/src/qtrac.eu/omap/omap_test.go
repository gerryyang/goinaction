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

// The tests here are very incomplete and just to show examples of how it
// can be done.
package omap_test

import (
    "qtrac.eu/omap"
    "strings"
    "testing"
)

func TestStringKeyOMapInsertion(t *testing.T) {
    wordForWord := omap.NewCaseFoldedKeyed()
    for _, word := range []string{"one", "Two", "THREE", "four", "Five"} {
        wordForWord.Insert(word, word)
    }
    var words []string
    wordForWord.Do(func(_, value interface{}) {
        words = append(words, value.(string))
    })
    actual, expected := strings.Join(words, ""), "FivefouroneTHREETwo"
    if actual != expected {
        t.Errorf("%q != %q", actual, expected)
    }
}

func TestIntKeyOMapFind(t *testing.T) {
    intMap := omap.NewIntKeyed()
    for _, number := range []int{9, 1, 8, 2, 7, 3, 6, 4, 5, 0} {
        intMap.Insert(number, number*10)
    }
    for _, number := range []int{0, 1, 5, 8, 9} {
        value, found := intMap.Find(number)
        if !found {
            t.Errorf("failed to find %d", number)
        }
        actual, expected := value.(int), number*10
        if actual != expected {
            t.Errorf("value is %d should be %d", actual, expected)
        }
    }
    for _, number := range []int{-1, -21, 10, 11, 148} {
        _, found := intMap.Find(number)
        if found {
            t.Errorf("should not have found %d", number)
        }
    }
}

func TestIntKeyOMapDelete(t *testing.T) {
    intMap := omap.NewIntKeyed()
    for _, number := range []int{9, 1, 8, 2, 7, 3, 6, 4, 5, 0} {
        intMap.Insert(number, number*10)
    }
    if intMap.Len() != 10 {
        t.Errorf("map len %d should be 10", intMap.Len())
    }
    length := 9
    for i, number := range []int{0, 1, 5, 8, 9} {
        if deleted := intMap.Delete(number); !deleted {
            t.Errorf("failed to delete %d", number)
        }
        if intMap.Len() != length-i {
            t.Errorf("map len %d should be %d", intMap.Len(), length-i)
        }
    }
    for _, number := range []int{-1, -21, 10, 11, 148} {
        if deleted := intMap.Delete(number); deleted {
            t.Errorf("should not have deleted nonexistent %d", number)
        }
    }
    if intMap.Len() != 5 {
        t.Errorf("map len %d should be 5", intMap.Len())
    }
}

func TestPassing(t *testing.T) {
    intMap := omap.NewIntKeyed()
    intMap.Insert(7, 7)
    passMap(intMap, t)
}

func passMap(m *omap.Map, t *testing.T) {
    for _, number := range []int{9, 3, 6, 4, 5, 0} {
        m.Insert(number, number)
    }
    if m.Len() != 7 {
        t.Errorf("should have %d items", 7)
    }
}

// Thanks to Russ Cox for improving these benchmarks
func BenchmarkOMapFindSuccess(b *testing.B) {
    b.StopTimer() // Don't time creation and population
    intMap := omap.NewIntKeyed()
    for i := 0; i < 1e6; i++ {
        intMap.Insert(i, i)
    }
    b.StartTimer() // Time the Find() method succeeding
    for i := 0; i < b.N; i++ {
        intMap.Find(i % 1e6)
    }
}

func BenchmarkOMapFindFailure(b *testing.B) {
    b.StopTimer() // Don't time creation and population
    intMap := omap.NewIntKeyed()
    for i := 0; i < 1e6; i++ {
        intMap.Insert(2*i, i)
    }
    b.StartTimer() // Time the Find() method failing
    for i := 0; i < b.N; i++ {
        intMap.Find(2*(i%1e6) + 1)
    }
}
