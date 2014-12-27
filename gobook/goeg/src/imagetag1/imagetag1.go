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

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU()) // Use all the machine's cores
    if len(os.Args) == 1 || os.Args[1] == "-h" || os.Args[1] == "--help" {
        fmt.Printf("usage: %s <image files>\n",
            filepath.Base(os.Args[0]))
        os.Exit(1)
    }

    files := commandLineFiles(os.Args[1:])
    for _, filename := range files {
        process(filename)
    }
}

func process(filename string) {
    if info, err := os.Stat(filename); err != nil ||
        (info.Mode()&os.ModeType != 0) {
        return // Ignore errors and nonregular files
    }
    file, err := os.Open(filename)
    if err != nil {
        return // Ignore errors
    }
    defer file.Close()
    config, _, err := image.DecodeConfig(file)
    if err != nil {
        return // Ignore errors
    }
    fmt.Printf(`<img src="%s" width="%d" height="%d" />`,
        filepath.Base(filename), config.Width, config.Height)
    fmt.Println()
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
