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
package fuzzybool

import "fmt"

type FuzzyBool struct{ value float32 }

func New(value interface{}) (*FuzzyBool, error) {
    amount, err := float32ForValue(value)
    return &FuzzyBool{amount}, err
}

func float32ForValue(value interface{}) (fuzzy float32, err error) {
    switch value := value.(type) { // shadow variable
    case float32:
        fuzzy = value
    case float64:
        fuzzy = float32(value)
    case int:
        fuzzy = float32(value)
    case bool:
        fuzzy = 0
        if value {
            fuzzy = 1
        }
    default:
        return 0, fmt.Errorf("float32ForValue(): %v is not a "+
            "number or Boolean", value)
    }
    if fuzzy < 0 {
        fuzzy = 0
    } else if fuzzy > 1 {
        fuzzy = 1
    }
    return fuzzy, nil
}

func (fuzzy *FuzzyBool) Set(value interface{}) (err error) {
    fuzzy.value, err = float32ForValue(value)
    return err
}

func (fuzzy *FuzzyBool) Copy() *FuzzyBool {
    return &FuzzyBool{fuzzy.value}
}

func (fuzzy *FuzzyBool) String() string {
    return fmt.Sprintf("%.0f%%", 100*fuzzy.value)
}

func (fuzzy *FuzzyBool) Not() *FuzzyBool {
    return &FuzzyBool{1 - fuzzy.value}
}

func (fuzzy *FuzzyBool) And(first *FuzzyBool,
    rest ...*FuzzyBool) *FuzzyBool {
    minimum := fuzzy.value
    rest = append(rest, first)
    for _, other := range rest {
        if minimum > other.value {
            minimum = other.value
        }
    }
    return &FuzzyBool{minimum}
}

func (fuzzy *FuzzyBool) Or(first *FuzzyBool,
    rest ...*FuzzyBool) *FuzzyBool {
    maximum := fuzzy.value
    rest = append(rest, first)
    for _, other := range rest {
        if maximum < other.value {
            maximum = other.value
        }
    }
    return &FuzzyBool{maximum}
}

func (fuzzy *FuzzyBool) Less(other *FuzzyBool) bool {
    return fuzzy.value < other.value
}

func (fuzzy *FuzzyBool) Equal(other *FuzzyBool) bool {
    return fuzzy.value == other.value
}

func (fuzzy *FuzzyBool) Bool() bool {
    return fuzzy.value >= .5
}

func (fuzzy *FuzzyBool) Float() float64 {
    return float64(fuzzy.value)
}
