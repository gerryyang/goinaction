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
    "bufio"
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "strings"
)

func main() {
    if len(os.Args) == 1 || os.Args[1] == "-h" || os.Args[1] == "--help" {
        fmt.Printf("usage: %s file\n", filepath.Base(os.Args[0]))
        os.Exit(1)
    }

    separators := []string{"\t", "*", "|", "•"}

    linesRead, lines := readUpToNLines(os.Args[1], 5)
    counts := createCounts(lines, separators, linesRead)
    separator := guessSep(counts, separators, linesRead)
    report(separator)
}

func createCounts(lines, separators []string, linesRead int) [][]int {
    counts := make([][]int, len(separators))
    for sepIndex := range separators {
        counts[sepIndex] = make([]int, linesRead)
        for lineIndex, line := range lines {
            counts[sepIndex][lineIndex] =
                strings.Count(line, separators[sepIndex])
        }
    }
    return counts
}

func guessSep(counts [][]int, separators []string, linesRead int) string {
    for sepIndex := range separators {
        same := true
        count := counts[sepIndex][0]
        for lineIndex := 1; lineIndex < linesRead; lineIndex++ {
            if counts[sepIndex][lineIndex] != count {
                same = false
                break
            }
        }
        if count > 0 && same {
            return separators[sepIndex]
        }
    }
    return ""
}

func report(separator string) {
    switch separator {
    case "":
        fmt.Println("whitespace-separated or not separated at all")
    case "\t":
        fmt.Println("tab-separated")
    default:
        fmt.Printf("%s-separated\n", separator)
    }
}

func readUpToNLines(filename string, maxLines int) (int, []string) {
    var file *os.File
    var err error
    if file, err = os.Open(filename); err != nil {
        log.Fatal("failed to open the file: ", err)
    }
    defer file.Close()
    lines := make([]string, maxLines)
    reader := bufio.NewReader(file)
    i := 0
    for ; i < maxLines; i++ {
        line, err := reader.ReadString('\n')
        if line != "" {
            lines[i] = line
        }
        if err != nil {
            if err == io.EOF {
                break
            }
            log.Fatal("failed to finish reading the file: ", err)
        }
    } // Return the subslice actually used; could be < maxLines
    return i, lines[:i]
}
