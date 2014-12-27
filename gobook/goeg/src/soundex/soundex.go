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
    "html"
    "io/ioutil"
    "log"
    "net/http"
    "sort"
    "strings"
)

const (
    pageTop    = `<!DOCTYPE HTML><html><head>
<style>.error{color:#FF0000;}</style></head>
<title>Soundex</title><body><h3>Soundex</h3>
<p>Compute soundex codes for a list of names.</p>`
    form       = `<form action="/" method="POST">
<label for="names">Names (comma or space-separated):</label><br />
<input type="text" name="names" size="30"><br />
<input type="submit" name="compute" value="Compute">
</form>`
    pageBottom = `</body></html>`
    error      = `<p class="error">%s</p>`
)

var digitForLetter = []rune{
    0, 1, 2, 3, 0, 1, 2, 0, 0, 2, 2, 4, 5,
    // A  B  C  D  E  F  G  H  I  J  K  L  M
    5, 0, 1, 2, 6, 2, 3, 0, 1, 0, 2, 0, 2}
    // N  O  P  Q  R  S  T  U  V  W  X  Y  Z

var testCases map[string]string

func main() {
    http.HandleFunc("/", homePage)
    var ok bool
    if testCases, ok = populateTestCases("soundex-test-data.txt"); ok {
        http.HandleFunc("/test", testPage)
    }
    if err := http.ListenAndServe(":9001", nil); err != nil {
        log.Fatal("failed to start server", err)
    }
}

func populateTestCases(filename string) (map[string]string, bool) {
    testCases := make(map[string]string)
    if lines, err := ioutil.ReadFile(filename); err != nil {
        log.Println(err)
        return testCases, false
    } else {
        for _, line := range strings.Split(string(lines), "\n") {
            if fields := strings.Fields(line); len(fields) == 2 {
                testCases[fields[1]] = fields[0]
            }
        }
    }
    return testCases, true
}

func homePage(writer http.ResponseWriter, request *http.Request) {
    err := request.ParseForm() // Must be called before writing response
    fmt.Fprint(writer, pageTop, form)
    if err != nil {
        fmt.Fprintf(writer, error, err)
    } else {
        if names := processRequest(request); len(names) > 0 {
            soundexes := make([]string, len(names))
            for i, name := range names {
                soundexes[i] = soundex(name)
            }
            fmt.Fprint(writer, formatResults(names, soundexes))
        }
    }
    fmt.Fprint(writer, pageBottom)
}

func processRequest(request *http.Request) (names []string) {
    if slice, found := request.Form["names"]; found && len(slice) > 0 {
        text := strings.Replace(slice[0], ",", " ", -1)
        names = strings.Fields(text)
    }
    return names
}

// "c - 'A'" produces a 0-based index, so 'A' -> 0, 'B' -> 1, etc.
// "'0' + digitForLetter[index]" converts a one digit integer into the
// equivalent Unicode character, i.e., 0 -> "0", 1 -> "1", etc.
func soundex(name string) string {
    name = strings.ToUpper(name)
    chars := []rune(name)
    var codes []rune
    codes = append(codes, chars[0])
    for i := 1; i < len(chars); i++ {
        char := chars[i]
        if i > 0 && char == chars[i-1] {
            continue
        }
        if index := char - 'A'; 0 <= index &&
            index < int32(len(digitForLetter)) &&
            digitForLetter[index] != 0 {
            codes = append(codes, '0'+digitForLetter[index])
        }
    }
    for len(codes) < 4 {
        codes = append(codes, '0')
    }
    return string(codes[:4])
}

func formatResults(names, soundexes []string) string {
    text := `<table border="1"><tr><th>Name</th><th>Soundex</th></tr>`
    for i := range names {
        text += "<tr><td>" + html.EscapeString(names[i]) + "</td><td>" +
            html.EscapeString(soundexes[i]) + "</td></tr>"
    }
    return text + "</table>"
}

func testPage(writer http.ResponseWriter, request *http.Request) {
    fmt.Fprint(writer, `<html><head><title>Soundex Test</title>
<style>.fail{color:#F00;} .pass{color:#0F0;}</style></head><body>
<table border="1"><tr><th>Name</th><th>Soundex</th>
<th>Expected</th><th>Test</th></tr>`)
    var names []string
    for name := range testCases {
        names = append(names, name)
    }
    sort.Strings(names)
    for _, name := range names {
        actual := soundex(name)
        expected := testCases[name]
        test := `<span class="fail">FAIL</span>`
        if actual == expected {
            test = `<span class="pass">PASS</span>`
        }
        fmt.Fprintf(writer, `<tr><td>%s</td><td>%s</td><td>%s</td>
<td>%s</td></tr>`, name, actual, expected, test)
    }
    fmt.Fprint(writer, "</table></body></html>")
}
