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
    "regexp"
    "runtime"
)

var workers = runtime.NumCPU()

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU()) // Use all the machine's cores
    if len(os.Args) != 2 || os.Args[1] == "-h" || os.Args[1] == "--help" {
        fmt.Printf("usage: %s <file.log>\n", filepath.Base(os.Args[0]))
        os.Exit(1)
    }
    filename := os.Args[1]
    info, err := os.Stat(filename)
    if err != nil {
        log.Fatal("failed to open the file:", err)
    }
    chunkSize := info.Size() / int64(workers)
    results := make(chan map[string]int, workers)
    getRx := regexp.MustCompile(`GET[ \t]+([^ \t\n]+[.]html?)`)
    for i := 0; i < workers; i++ {
        offset := int64(i) * chunkSize
        if i + 1 == workers {
            chunkSize *= 2 // Make sure we read all of the last chunk
        }
        go processLines(results, getRx, filename, offset, chunkSize)
    }
    totalForPage := make(map[string]int)
    merge(results, totalForPage)
    showResults(totalForPage)
}

func processLines(results chan<- map[string]int, getRx *regexp.Regexp,
    filename string, offset, chunkSize int64) {
    file, err := os.Open(filename)
    if err != nil {
        log.Fatal("failed to open the file:", err)
    }
    defer file.Close()
    file.Seek(offset, 0)
    var bytesRead int64
    reader := bufio.NewReader(file)
    if (offset > 0) { // Find first whole line
        line, err := reader.ReadString('\n')
        if err != nil {
            log.Fatal("failed to read the file:", err)
        }
        bytesRead = int64(len(line))
    }
    countForPage := make(map[string]int)
    for bytesRead < chunkSize {
        line, err := reader.ReadString('\n')
        if line != "" {
            bytesRead += int64(len(line))
            if matches := getRx.FindStringSubmatch(line); matches != nil {
                countForPage[matches[1]]++
            }
        }
        if err != nil {
            if err != io.EOF {
                log.Println("failed to finish reading the file:", err)
            }
            break
        }
    }
    results <- countForPage
}

func merge(results <-chan map[string]int, totalForPage map[string]int) {
    for i := 0; i < workers; i++ {
        countForPage := <-results
        for page, count := range countForPage {
            totalForPage[page] += count
        }
    }
}

func showResults(totalForPage map[string]int) {
    for page, count := range totalForPage {
        fmt.Printf("%8d %s\n", count, page)
    }
}
