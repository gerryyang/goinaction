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
    "image"
    "image/color"
    "shaper_ans1/shapes"
)

func main() {
    img := shapes.FilledImage(420, 220, image.White)
    fill := color.RGBA{200, 200, 200, 0xFF} // light gray
    for i := 0; i < 10; i++ {
        width, height := 40+(20*i), 20+(10*i)
        rectangle := shapes.NewRectangle(fill,
            image.Rect(0, 0, width, height))
        rectangle.SetFilled(true)
        x := 10 + (20 * i)
        for j := i / 2; j >= 0; j-- {
            rectangle.Draw(img, x+j, (x/2)+j)
        }
        fill.R -= uint8(i * 5)
        fill.G = fill.R
        fill.B = fill.R
    }
    shapes.SaveImage(img, "rectangle.png")
}
