// Copyright © 2011-13 Qtrac Ltd.
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
    "net/http"
    "sort"
    "strconv"
    "strings"
)

const (
    pageTop    = `<!DOCTYPE HTML><html><head>
<style>.error{color:#FF0000;}</style></head><title>Statistics</title>
<body><h3>Statistics</h3>
<p>Computes basic statistics for a given list of numbers</p>`
    form       = `<form action="/" method="POST">
<label for="numbers">Numbers (comma or space-separated):</label><br />
<input type="text" name="numbers" size="30"><br />
<input type="submit" name="calculate" value="Calculate">
<input type="submit" name="panic" value="Panic!">
</form>`
    pageBottom = `</body></html>`
    anError    = `<p class="error">%s</p>`
)

type statistics struct {
    numbers []float64
    mean    float64
    median  float64
    modes   []float64
    stdDev  float64
}

func main() {
    log.SetFlags(0) // Don't log timestamps
    http.HandleFunc("/", homePage)
    if err := http.ListenAndServe(":9001", nil); err != nil {
        log.Fatal("failed to start server", err)
    }
}

func homePage(writer http.ResponseWriter, request *http.Request) {
    defer func() { // Needed for every page
        if x := recover(); x != nil {
            log.Printf("[%v] caught panic: %v", request.RemoteAddr, x)
        }
    }()

    err := request.ParseForm() // Must be called before writing response
    fmt.Fprint(writer, pageTop, form)
    if err != nil {
        fmt.Fprintf(writer, anError, err)
    } else {
        if numbers, message, ok := processRequest(request); ok {
            stats := getStats(numbers)
            fmt.Fprint(writer, formatStats(stats))
            log.Printf("[%v] served OK", request.RemoteAddr)
        } else if message != "" {
            fmt.Fprintf(writer, anError, message)
            log.Printf("[%v] bad request: %v", request.RemoteAddr,
                message)
        }
    }
    fmt.Fprint(writer, pageBottom)
}

func processRequest(request *http.Request) ([]float64, string, bool) {
    if _, found := request.Form["panic"]; found { // Fake a panic
        panic("user clicked panic button!")
    }
    var numbers []float64
    if slice, found := request.Form["numbers"]; found && len(slice) > 0 {
        text := strings.Replace(slice[0], ",", " ", -1)
        for _, field := range strings.Fields(text) {
            if x, err := strconv.ParseFloat(field, 64); err != nil {
                return numbers, "'" + field + "' is invalid", false
            } else {
                numbers = append(numbers, x)
            }
        }
    }
    if len(numbers) == 0 {
        return numbers, "", false // no data first time form is shown
    }
    return numbers, "", true
}

func formatStats(stats statistics) string {
    return fmt.Sprintf(`<table border="1">
<tr><th colspan="2">Results</th></tr>
<tr><td>Numbers</td><td>%v</td></tr>
<tr><td>Count</td><td>%d</td></tr>
<tr><td>Mean</td><td>%f</td></tr>
<tr><td>Median</td><td>%f</td></tr>
<tr><td>Mode</td><td>%v</td></tr>
<tr><td>Std. Dev.</td><td>%f</td></tr>
</table>`, stats.numbers, len(stats.numbers), stats.mean, stats.median,
        stats.modes, stats.stdDev)
}

func getStats(numbers []float64) (stats statistics) {
    stats.numbers = numbers
    sort.Float64s(stats.numbers)
    stats.mean = sum(numbers) / float64(len(numbers))
    stats.median = median(numbers)
    stats.modes = mode(numbers)
    stats.stdDev = stdDev(numbers, stats.mean)
    return stats
}

func sum(numbers []float64) (total float64) {
    for _, x := range numbers {
        total += x
    }
    return total
}

func median(numbers []float64) float64 {
    middle := len(numbers) / 2
    result := numbers[middle]
    if len(numbers)%2 == 0 {
        result = (result + numbers[middle-1]) / 2
    }
    return result
}

func mode(numbers []float64) (modes []float64) {
    frequencies := make(map[float64]int, len(numbers))
    highestFrequency := 0
    for _, x := range numbers {
        frequencies[x]++
        if frequencies[x] > highestFrequency {
            highestFrequency = frequencies[x]
        }
    }
    for x, frequency := range frequencies {
        if frequency == highestFrequency {
            modes = append(modes, x)
        }
    }
    if highestFrequency == 1 || len(modes) == len(frequencies) {
        modes = modes[:0] // Or: modes = []float64{}
    }
    sort.Float64s(modes)
    return modes
}

func stdDev(numbers []float64, mean float64) float64 {
    total := 0.0
    for _, number := range numbers {
        total += math.Pow(number-mean, 2)
    }
    variance := total / float64(len(numbers)-1)
    return math.Sqrt(variance)
}
