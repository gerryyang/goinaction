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
    "fmt"
    "image"
    "os"
    "path/filepath"
    "runtime"

    _ "image/gif"
    _ "image/jpeg"
    _ "image/png"
)

var workers = runtime.NumCPU()

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU()) // Use all the machine's cores
    if len(os.Args) == 1 || os.Args[1] == "-h" || os.Args[1] == "--help" {
        fmt.Printf("usage: %s <image files>\n",
            filepath.Base(os.Args[0]))
        os.Exit(1)
    }

    files := commandLineFiles(os.Args[1:])
    jobs := make(chan string, workers*16)
    results := make(chan string)
    done := make(chan struct{}, workers)

    go addJobs(files, jobs)
    for i := 0; i < workers; i++ {
        go doJobs(done, results, jobs)
    }
    waitAndProcessResults(done, results)
}

func addJobs(files []string, jobs chan<- string) {
    for _, filename := range files {
        jobs <- filename
    }
    close(jobs)
}

func doJobs(done chan<- struct{}, results chan<- string,
    jobs <-chan string) {
    for job := range jobs {
        if result, ok := process(job); ok {
            results <- result
        }
    }
    done <- struct{}{}
}

func waitAndProcessResults(done <-chan struct{}, results <-chan string) {
    for working := workers; working > 0; {
        select { // Blocking
        case result := <-results:
            fmt.Println(result)
        case <-done:
            working--
        }
    }
DONE:
    for {
        select { // Nonblocking
        case result := <-results:
            fmt.Println(result)
        default:
            break DONE
        }
    }
}

func process(filename string) (string, bool) {
    if info, err := os.Stat(filename); err != nil ||
        (info.Mode()&os.ModeType == 1) {
        return "", false // Ignore errors and nonregular files
    }
    file, err := os.Open(filename)
    if err != nil {
        return "", false // Ignore errors
    }
    defer file.Close()
    config, _, err := image.DecodeConfig(file)
    if err != nil {
        return "", false // Ignore errors
    }
    return fmt.Sprintf(`<img src="%s" width="%d" height="%d" />`,
        filepath.Base(filename), config.Width, config.Height), true
}

func commandLineFiles(files []string) []string {
    if runtime.GOOS == "windows" {
        args := make([]string, 0, len(files))
        for _, name := range files {
            if matches, err := filepath.Glob(name); err != nil {
                args = append(args, name) // Invalid pattern
            } else if matches != nil { // At least one match
                args = append(args, matches...)
            }
        }
        return args
    }
    return files
}
