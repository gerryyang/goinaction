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
    "bytes"
    "fmt"
    "strings"
)

type memoizeFunction func(int, ...int) interface{}

var Fibonacci memoizeFunction
var RomanForDecimal memoizeFunction

func init() {
    // Of course the iterative version of fibonacci (that doesn't need
    // memoize) is much more efficient.
    Fibonacci = Memoize(func(x int, xs ...int) interface{} {
        if x < 2 {
            return x
        }
        return Fibonacci(x-1).(int) + Fibonacci(x-2).(int)
    })

    decimals := []int{1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1}
    romans := []string{"M", "CM", "D", "CD", "C", "XC", "L", "XL", "X",
        "IX", "V", "IV", "I"}
    RomanForDecimal = Memoize(func(x int, xs ...int) interface{} {
        if x < 0 || x > 3999 {
            panic("RomanForDecimal() only handles integers [0, 3999]")
        }
        var buffer bytes.Buffer
        for i, decimal := range decimals {
            remainder := x / decimal
            x %= decimal
            if remainder > 0 {
                buffer.WriteString(strings.Repeat(romans[i], remainder))
            }
        }
        return buffer.String()
    })
}

func main() {
    fmt.Println("Fibonacci(45) =", Fibonacci(45).(int))
    for _, x := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13,
        14, 15, 16, 17, 18, 19, 20, 25, 30, 40, 50, 60, 69, 70, 80,
        90, 99, 100, 200, 300, 400, 500, 600, 666, 700, 800, 900,
        1000, 1009, 1444, 1666, 1945, 1997, 1999, 2000, 2008, 2010,
        2012, 2500, 3000, 3999} {
        fmt.Printf("%4d = %s\n", x, RomanForDecimal(x).(string))
    }
}

func Memoize(function memoizeFunction) memoizeFunction {
    cache := make(map[string]interface{})
    return func(x int, xs ...int) interface{} {
        key := fmt.Sprint(x)
        for _, i := range xs {
            key += fmt.Sprintf(",%d", i)
        }
        if value, found := cache[key]; found {
            return value
        }
        value := function(x, xs...)
        cache[key] = value
        return value
    }
}
