// Copyright © 2010-12 Qtrac Ltd.
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
	"stacker/stack"
	"strings"
)

func main() {
	var haystack stack.Stack
	haystack.Push("hay")
	haystack.Push(-15)
	haystack.Push([]string{"pin", "clip", "needle"})
	haystack.Push(81.52)
	for {
		item, err := haystack.Pop()
		if err != nil {
			break
		}
		fmt.Println(item)
	}

	var aStack stack.Stack
	aStack.Push("Aarvark")
	aStack.Push(5)
	aStack.Push(19)
	x, err := aStack.Top()
	fmt.Println(x)
	aStack.Push(-6e-4)
	aStack.Push("Baker")
	aStack.Push(-3)
	aStack.Push("Cake")
	aStack.Push("Dancer")
	x, err = aStack.Top()
	fmt.Println(x)
	aStack.Push(11.7)
	fmt.Println("stack is empty", aStack.IsEmpty())
	fmt.Printf("Len() == %d  Cap == %d\n", aStack.Len(), aStack.Cap())
	difference := aStack.Cap() - aStack.Len()
	for i := 0; i < difference; i++ {
		aStack.Push(strings.Repeat("*", difference-i))
	}
	fmt.Printf("Len() == %d  Cap == %d\n", aStack.Len(), aStack.Cap())
	for aStack.Len() > 0 {
		x, _ = aStack.Pop()
		fmt.Printf("%T %v\n", x, x)
	}
	fmt.Println("stack is empty", aStack.IsEmpty())
	x, err = aStack.Pop()
	fmt.Println(x, err)
	x, err = aStack.Top()
	fmt.Println(x, err)
}
