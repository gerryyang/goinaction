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
    "os"
    "path/filepath"
    "unicode/utf8"
)

var IsPalindrome func(string) bool // Holds a reference to a function


func init() {
    if len(os.Args) > 1 &&
        (os.Args[1] == "-a" || os.Args[1] == "--ascii") {
        os.Args = append(os.Args[:1], os.Args[2:]...) // Strip out arg.
        IsPalindrome = func(s string) bool {  // Simple ASCII-only version
            if len(s) <= 1 {
                return true
            }
            if s[0] != s[len(s)-1] {
                return false
            }
            return IsPalindrome(s[1 : len(s)-1])
        }
    } else {
        IsPalindrome = func(s string) bool { // UTF-8 version
            if utf8.RuneCountInString(s) <= 1 {
                return true
            }
            first, sizeOfFirst := utf8.DecodeRuneInString(s)
            last, sizeOfLast := utf8.DecodeLastRuneInString(s)
            if first != last {
                return false
            }
            return IsPalindrome(s[sizeOfFirst : len(s)-sizeOfLast])
        }
    }
}

func main() {
    if len(os.Args) == 1 {
        fmt.Printf("usage: %s [-a|--ascii] word1 [word2 [... wordN]]\n",
            filepath.Base(os.Args[0]))
        os.Exit(1)
    }
    words := os.Args[1:]
    for _, word := range words {
        fmt.Printf("%5t %q\n", IsPalindrome(word), word)
    }
}
