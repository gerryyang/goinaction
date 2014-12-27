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

package linkutil_test

import (
    "bufio"
    "io"
    "linkcheck/linkutil"
    "os"
    "reflect"
    "sort"
    "testing"
)

func TestLinksFromReader(t *testing.T) {
    file, err := os.Open("index.html")
    if err != nil {
        t.Fatal(err)
    }
    defer file.Close()
    links, err := linkutil.LinksFromReader(file)
    if err != nil {
        t.Fatal(err)
    }
    sort.Strings(links)
    file, err = os.Open("index.links")
    if err != nil {
        t.Fatal(err)
    }
    defer file.Close()
    reader := bufio.NewReader(file)
    var lines []string
    for {
        line, err := reader.ReadString('\n')
        if line != "" {
            lines = append(lines, line[:len(line)-1])
        }
        if err != nil {
            if err != io.EOF {
                t.Fatal(err)
            }
            break
        }
    }
    sort.Strings(lines)
    if !reflect.DeepEqual(links, lines) {
        for i := 0; i < len(links); i++ {
            if i < len(lines) {
                if links[i] != lines[i] {
                    t.Fatalf("%q != %q", links[i], lines[i])
                }
            } else {
                t.Fatalf("found more links than lines, e.g.: %q", links[i])
            }
        }
        t.Fatalf("found fewer links than lines (%d vs. %d)", len(links),
            len(lines))
    }
}
