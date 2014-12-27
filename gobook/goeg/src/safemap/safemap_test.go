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

package safemap_test

import (
    "fmt"
    "safemap"
    "sync"
    "testing"
)

func TestSafeMap(t *testing.T) {
    store := safemap.New()
    fmt.Printf("Initially has %d items\n", store.Len())

    deleted := []int{0, 2, 3, 5, 7, 20, 399, 25, 30, 1000, 91, 97, 98, 99}

    var waiter sync.WaitGroup

    waiter.Add(1)
    go func() { // Concurrent Inserter
        for i := 0; i < 100; i++ {
            store.Insert(fmt.Sprintf("0x%04X", i), i)
            if i > 0 && i%15 == 0 {
                fmt.Printf("Inserted %d items\n", store.Len())
            }
        }
        fmt.Printf("Inserted %d items\n", store.Len())
        waiter.Done()
    }()

    waiter.Add(1)
    go func() { // Concurrent Deleter
        for _, i := range deleted {
            key := fmt.Sprintf("0x%04X", i)
            before := store.Len()
            store.Delete(key)
            fmt.Printf("Deleted m[%s] (%d) before=%d after=%d\n",
                key, i, before, store.Len())
        }
        waiter.Done()
    }()

    waiter.Add(1)
    go func() { // Concurrent Finder
        for _, i := range deleted {
            for _, j := range []int{i, i + 1} {
                key := fmt.Sprintf("0x%04X", j)
                value, found := store.Find(key)
                if found {
                    fmt.Printf("Found m[%s] == %d\n", key, value)
                } else {
                    fmt.Printf("Not found m[%s] (%d)\n", key, j)
                }
            }
        }
        waiter.Done()
    }()

    waiter.Wait()

    updater := func(value interface{}, found bool) interface{} {
        if found {
            return value.(int) * 1000
        }
        return 1
    }
    for _, i := range []int{5, 10, 15, 20, 25, 30, 35} {
        key := fmt.Sprintf("0x%04X", i)
        if value, found := store.Find(key); found {
            fmt.Printf("Original m[%s] == %d\t", key, value)
            store.Update(key, updater)
            if value, found := store.Find(key); found {
                fmt.Printf("Updated m[%s] == %5d\n", key, value)
            }
        }
    }

    fmt.Printf("Finished with %d items\n", store.Len())
    // not needed here but useful if you want to free up the goroutine
    data := store.Close()
    fmt.Println("Closed")
    fmt.Printf("len == %d\n", len(data))
    //for k, v := range data { fmt.Printf("%s = %v\n", k, v) }
}
