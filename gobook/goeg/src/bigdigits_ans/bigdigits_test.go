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
    "bytes"
    "io/ioutil"
    "log"
    "os"
    "os/exec"
    "path/filepath"
    "testing"
)

func TestBigDigits(t *testing.T) {
    log.SetFlags(0)
    log.Println("TEST bigdigits_ans")

    path, _ := os.Getwd()
    expected, err := ioutil.ReadFile(filepath.Join(path, "0123456789.txt"))
    if err != nil {
        t.Fatal(err)
    }
    executable := filepath.Join(path, "bigdigits_ans")
    reader, writer, err := os.Pipe()
    if err != nil {
        t.Fatal(err)
    }
    command := exec.Command(executable, "-b", "0123456789")
    command.Stdout = writer
    err = command.Run()
    if err != nil {
        t.Fatal(err)
    }
    writer.Close()
    actual, err := ioutil.ReadAll(reader)
    if err != nil {
        t.Fatal(err)
    }
    reader.Close()
    if bytes.Compare(actual, expected) != 0 {
        t.Fatal("actual != expected")
    }
}
