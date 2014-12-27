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

// Algorithm taken from:
// http://en.literateprograms.org/Pi_with_Machin%27s_formula_%28Python%29

package main

import (
    "fmt"
    "math/big"
    "os"
    "path/filepath"
    "strconv"
)

func main() {
    places := handleCommandLine(1000)
    scaledPi := fmt.Sprint(π(places))
    fmt.Printf("3.%s\n", scaledPi[1:])
}

func handleCommandLine(defaultValue int) int {
    if len(os.Args) > 1 {
        if os.Args[1] == "-h" || os.Args[1] == "--help" {
            usage := "usage: %s [digits]\n e.g.: %s 10000"
            app := filepath.Base(os.Args[0])
            fmt.Fprintln(os.Stderr, fmt.Sprintf(usage, app, app))
            os.Exit(1)
        }
        if x, err := strconv.Atoi(os.Args[1]); err != nil {
            fmt.Fprintf(os.Stderr, "ignoring invalid number of "+
                "digits: will display %d\n", defaultValue)
        } else {
            return x
        }
    }
    return defaultValue
}

func π(places int) *big.Int {
    digits := big.NewInt(int64(places))
    unity := big.NewInt(0)
    ten := big.NewInt(10)
    exponent := big.NewInt(0)
    unity.Exp(ten, exponent.Add(digits, ten), nil)
    pi := big.NewInt(4)
    left := arccot(big.NewInt(5), unity)
    left.Mul(left, big.NewInt(4))
    right := arccot(big.NewInt(239), unity)
    left.Sub(left, right)
    pi.Mul(pi, left)
    return pi.Div(pi, big.NewInt(0).Exp(ten, ten, nil))
}

func arccot(x, unity *big.Int) *big.Int {
    sum := big.NewInt(0)
    sum.Div(unity, x)
    xpower := big.NewInt(0)
    xpower.Div(unity, x)
    n := big.NewInt(3)
    sign := big.NewInt(-1)
    zero := big.NewInt(0)
    square := big.NewInt(0)
    square.Mul(x, x)
    for {
        xpower.Div(xpower, square)
        term := big.NewInt(0)
        term.Div(xpower, n)
        if term.Cmp(zero) == 0 {
            break
        }
        addend := big.NewInt(0)
        sum.Add(sum, addend.Mul(sign, term))
        sign.Neg(sign)
        n.Add(n, big.NewInt(2))
    }
    return sum
}
