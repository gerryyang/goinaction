// Copyright © 2010-12 Qtrac Ltd.
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
    "log"
    "math"
    "math/cmplx"
    "net/http"
    "strconv"
)

const (
    decimals     = 3
    pageTop      = `<!DOCTYPE HTML><html><head>
<style>.error{color:#FF0000;}</style></head>
<title>Quadratic Equation Solver</title><body>
<h3>Quadratic Equation Solver</h3><p>Solves equations of the form
a<i>x</i>² + b<i>x</i> + c</p>`
    form         = `<form action="/" method="POST">
<input type="text" name="a" size="1"><label for="a"><i>x</i>²</label> +
<input type="text" name="b" size="1"><label for="b"><i>x</i></label> +
<input type="text" name="c" size="1"><label for="c"> →</label>
<input type="submit" name="calculate" value="Calculate">
</form>`
    pageBottom   = "</body></html>"
    error        = `<p class="error">%s</p>`
    solution     = "<p>%s → %s</p>"
    oneSolution  = "<i>x</i>=%s"
    twoSolutions = "<i>x</i>=%s or <i>x</i>=%s"
    noSolution   = "<i>there are no solutions</i>"
)

func main() {
    http.HandleFunc("/", homePage)
    if err := http.ListenAndServe(":9001", nil); err != nil {
        log.Fatal("failed to start server", err)
    }
}

func homePage(writer http.ResponseWriter, request *http.Request) {
    err := request.ParseForm() // Must be called before writing response
    fmt.Fprint(writer, pageTop, form)
    if err != nil {
        fmt.Fprintf(writer, error, err)
    } else {
        if floats, message, ok := processRequest(request); ok {
            question := formatQuestion(request.Form)
            x1, x2 := solve(floats)
            answer := formatSolutions(x1, x2)
            fmt.Fprintf(writer, solution, question, answer)
        } else if message != "" {
            fmt.Fprintf(writer, error, message)
        }
    }
    fmt.Fprint(writer, pageBottom)
}

func processRequest(request *http.Request) ([3]float64, string, bool) {
    var floats [3]float64
    count := 0
    for index, key := range []string{"a", "b", "c"} {
        if slice, found := request.Form[key]; found && len(slice) > 0 {
            if slice[0] != "" {
                if x, err := strconv.ParseFloat(slice[0], 64);
                    err != nil {
                    return floats, "'" + slice[0] + "' is invalid", false
                } else {
                    floats[index] = x
                }
            } else { // as a courtesy to users treat blanks as 0
                request.Form[key][0] = "0"
                floats[index] = 0
            }
            count++
        }
    }
    if count != 3 { // the first time the form is empty;
        return floats, "", false // this isn't an error but there's
    } // nothing to calculate
    if EqualFloat(floats[0], 0, -1) {
        return floats, "the x² factor may not be 0", false
    }
    return floats, "", true
}

func formatQuestion(form map[string][]string) string {
    result := formatSignAndNumber("", form["a"][0], "<i>x</i>²")
    result += formatSignAndNumber(" ", form["b"][0], "<i>x</i>")
    result += formatSignAndNumber(" ", form["c"][0], "")
    return result
}

func formatSignAndNumber(signPad, number, suffix string) string {
    if number == "" || number == "0" || number == "0.0" {
        return ""
    }
    var sign string
    if signPad != "" {
        sign = signPad + "+" + signPad
    }
    if number[0] == '-' {
        sign = signPad + "-" + signPad
        number = number[1:]
    }
    if suffix != "" && number == "1" {
        return sign + suffix
    }
    return sign + number + suffix
}

func formatSolutions(x1, x2 complex128) string {
    exactlyOneSolution := false
    if cmplx.IsNaN(x1) && cmplx.IsNaN(x2) {
        return noSolution
    }
    if cmplx.IsNaN(x1) {
        exactlyOneSolution = true
        x1 = x2
    } else if cmplx.IsNaN(x2) || EqualComplex(x1, x2) {
        exactlyOneSolution = true
    }
    if exactlyOneSolution {
        return fmt.Sprintf(oneSolution, formatComplex(x1))
    }
    return fmt.Sprintf(twoSolutions, formatComplex(x1), formatComplex(x2))
}

func formatComplex(x complex128) string {
    if EqualFloat(imag(x), 0, -1) {
        return fmt.Sprintf("%.*f", decimals, real(x))
    }
    return fmt.Sprintf("%.*f", decimals, x)
}

func solve(floats [3]float64) (complex128, complex128) {
    a, b, c := complex(floats[0], 0), complex(floats[1], 0),
        complex(floats[2], 0)
    root := cmplx.Sqrt(cmplx.Pow(b, 2) - (4 * a * c))
    x1 := (-b + root) / (2 * a)
    x2 := (-b - root) / (2 * a)
    return x1, x2
}

// EqualFloat() returns true if x and y are approximately equal to the
// given limit. Pass a limit of -1 to get the greatest accuracy the machine
// can manage.
func EqualFloat(x, y, limit float64) bool {
    if limit <= 0.0 {
        limit = math.SmallestNonzeroFloat64
    }
    return math.Abs(x-y) <=
        (limit * math.Min(math.Abs(x), math.Abs(y)))
}

func EqualComplex(x, y complex128) bool {
    return EqualFloat(real(x), real(y), -1) &&
        EqualFloat(imag(x), imag(y), -1)
}
