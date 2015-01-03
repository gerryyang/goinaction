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

package stack_test

import (
    "stacker/stack"
    "testing"
)

func TestStack(t *testing.T) {
    count := 1
    var aStack stack.Stack
    assertTrue(t, aStack.Len() == 0, "expected empty Stack", count) // 1
    count++
    assertTrue(t, aStack.Cap() == 0, "expected empty Stack", count) // 2
    count++
    assertTrue(t, aStack.IsEmpty(), "expected empty Stack", count) // 3
    count++
    value, err := aStack.Pop()
    assertTrue(t, value == nil, "expected nil value", count) // 4
    count++
    assertTrue(t, err != nil, "expected error", count) // 5
    count++
    value1, err := aStack.Top()
    assertTrue(t, value1 == nil, "expected nil value", count) // 6
    count++
    assertTrue(t, err != nil, "expected error", count) // 7
    count++
    aStack.Push(1)
    aStack.Push(2)
    aStack.Push("three")
    assertTrue(t, aStack.Len() == 3, "expected nonempty Stack", count) // 8
    count++
    assertTrue(t, aStack.IsEmpty() == false, "expected nonempty Stack",
        count) // 9
    count++
    value2, err := aStack.Pop()
    assertEqualString(t, value2.(string), "three", "unexpected text",
        count) // 10
    count++
    assertTrue(t, err == nil, "no error expected", count) // 11
    count++
    value3, err := aStack.Top()
    assertTrue(t, value3 == 2, "unexpected number", count) // 12
    count++
    assertTrue(t, err == nil, "no error expected", count) // 13
    count++
    aStack.Pop()
    assertTrue(t, aStack.Len() == 1, "expected nonempty Stack", count) //14
    count++
    assertTrue(t, aStack.IsEmpty() == false, "expected nonempty Stack",
        count) // 15
    count++
    value4, err := aStack.Pop()
    assertTrue(t, value4 == 1, "unexpected number", count) // 16
    count++
    assertTrue(t, err == nil, "no error expected", count) // 17
    count++
    assertTrue(t, aStack.Len() == 0, "expected empty Stack", count) // 18
    count++
    assertTrue(t, aStack.IsEmpty(), "expected empty Stack", count) // 19
    count++
}

// assertTrue() calls testing.T.Error() with the given message if the
// condition is false.
func assertTrue(t *testing.T, condition bool, message string, id int) {
    if !condition {
        t.Errorf("#%d: %s", id, message)
    }
}

// assertEqualString() calls testing.T.Error() with the given message if
// the given strings are not equal.
func assertEqualString(t *testing.T, a, b string, message string, id int) {
    if a != b {
        t.Errorf("#%d: %s \"%s\" !=\n\"%s\"", id, message, a, b)
    }
}
