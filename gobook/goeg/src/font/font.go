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

package font

import (
    "fmt"
    "log"
    "unicode/utf8"
)

type Font struct {
    family string
    size   int
}

func New(family string, size int) *Font {
    return &Font{saneFamily("sans-serif", family), saneSize(10, size)}
}

func (font *Font) Family() string { return font.family }

func (font *Font) SetFamily(family string) {
    font.family = saneFamily(font.family, family)
}

func (font *Font) Size() int { return font.size }

func (font *Font) SetSize(size int) {
    font.size = saneSize(font.size, size)
}

func (font *Font) String() string {
    return fmt.Sprintf("{font-family: %q; font-size: %dpt;}", font.family,
        font.size)
}

func saneFamily(oldFamily, newFamily string) string {
    if len(newFamily) < utf8.UTFMax &&
        utf8.RuneCountInString(newFamily) < 1 {
        log.Printf("font.saneFamily(): ignored invalid family '%s'",
            newFamily)
        return oldFamily
    }
    return newFamily
}

func saneSize(oldSize, newSize int) int {
    if newSize < 5 || newSize > 144 {
        log.Printf("font.saneSize(): ignored invalid size '%d'", newSize)
        return oldSize
    }
    return newSize
}
