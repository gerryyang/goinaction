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

package font_test

import (
    "fmt"
    "font"
    "testing"
)

func TestFont(t *testing.T) {
    bodyFont := font.New("Nimbus Sans", 10)
    titleFont := font.New("serif", 11)
    f1(bodyFont, titleFont, t)
}

func f1(bodyFont, titleFont *font.Font, t *testing.T) {
    if bodyFont.String() !=
        `{font-family: "Nimbus Sans"; font-size: 10pt;}` {
        t.Fatal("#1 bodyFont invalid CSS")
    }
    if bodyFont.Size() != 10 || bodyFont.Family() != "Nimbus Sans" {
        t.Fatal("#2 bodyFont invalid attributes")
    }
    bodyFont.SetSize(12)
    if bodyFont.Size() != 12 || bodyFont.Family() != "Nimbus Sans" {
        t.Fatal("#3 bodyFont invalid attributes")
    }
    if bodyFont.String() !=
        `{font-family: "Nimbus Sans"; font-size: 12pt;}` {
        t.Fatal("#4 bodyFont invalid CSS")
    }
    bodyFont.SetFamily("")
    if bodyFont.Size() != 12 || bodyFont.Family() != "Nimbus Sans" {
        t.Fatal("#5 bodyFont invalid attributes")
    }

    if titleFont.String() != `{font-family: "serif"; font-size: 11pt;}` {
        t.Fatal("#6 titleFont invalid CSS")
    }
    if titleFont.Size() != 11 || titleFont.Family() != "serif" {
        t.Fatal("#7 titleFont invalid attributes")
    }
    titleFont.SetFamily("Helvetica")
    titleFont.SetSize(20)
    if titleFont.Size() != 20 || titleFont.Family() != "Helvetica" {
        t.Fatal("#8 titleFont invalid attributes")
    }

    f2(bodyFont, titleFont)
}

func f2(bodyFont, titleFont *font.Font) {
    fmt.Println(bodyFont)
    fmt.Println(titleFont)
}
