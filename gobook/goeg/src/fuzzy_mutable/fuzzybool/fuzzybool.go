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

type FuzzyBool float32

func New(value interface{}) (*FuzzyBool, error) {
    var fuzzy FuzzyBool
    result, err := float32ForValue(value)
    fuzzy = FuzzyBool(result)
    return &fuzzy, err
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

func (fuzzy *FuzzyBool) Set(value interface{}) error {
    result, err := float32ForValue(value)
    *fuzzy = FuzzyBool(result)
    return err
}

func (fuzzy *FuzzyBool) Copy() *FuzzyBool {
    result := FuzzyBool(*fuzzy)
    return &result
}

func (fuzzy *FuzzyBool) String() string {
    return fmt.Sprintf("%.0f%%", 100*float32(*fuzzy))
}

func (fuzzy *FuzzyBool) Not() *FuzzyBool {
    result := FuzzyBool(1 - float32(*fuzzy))
    return &result
}

func (fuzzy *FuzzyBool) And(first *FuzzyBool,
    rest ...*FuzzyBool) *FuzzyBool {
    minimum := *fuzzy
    rest = append(rest, first)
    for _, other := range rest {
        if minimum > *other {
            minimum = *other
        }
    }
    result := FuzzyBool(minimum)
    return &result
}

func (fuzzy *FuzzyBool) Or(first *FuzzyBool,
    rest ...*FuzzyBool) *FuzzyBool {
    maximum := *fuzzy
    rest = append(rest, first)
    for _, other := range rest {
        if maximum < *other {
            maximum = *other
        }
    }
    result := FuzzyBool(maximum)
    return &result
}

func (fuzzy *FuzzyBool) Less(other *FuzzyBool) bool {
    return *fuzzy < *other
}

func (fuzzy *FuzzyBool) Equal(other *FuzzyBool) bool {
    return *fuzzy == *other
}

func (fuzzy *FuzzyBool) Bool() bool {
    return float32(*fuzzy) >= .5
}

func (fuzzy *FuzzyBool) Float() float64 {
    return float64(*fuzzy)
}
