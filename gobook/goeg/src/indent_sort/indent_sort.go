// Copyright © 2011-14 Qtrac Ltd.
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
    "sort"
    "strings"
)

var original = []string{
    "Nonmetals",
    "    Hydrogen",
    "    Carbon",
    "    Nitrogen",
    "    Oxygen",
    "Inner Transitionals",
    "    Lanthanides",
    "        Europium",
    "        Cerium",
    "    Actinides",
    "        Uranium",
    "        Plutonium",
    "        Curium",
    "Alkali Metals",
    "    Lithium",
    "    Sodium",
    "    Potassium",
}

func main() {
    fmt.Println("|     Original      |       Sorted      |")
    fmt.Println("|-------------------|-------------------|")
    sorted := SortedIndentedStrings(original) // original is a []string
    for i := range original {                 // set in a global var
        fmt.Printf("|%-19s|%-19s|\n", original[i], sorted[i])
    }
}

/*
   Given a []string that has items with different levels of indent that
   are used to indicate parent → child relationships, sorts the items
   case-insensitively with child items sorted underneath their parent
   items, and so on recursively to any level of depth.
   The amount of indent per level is computed by finding the first
   indented item. Indentation must either be one or more spaces or one or
   more tabs.
*/
func SortedIndentedStrings(slice []string) []string {
    entries := populateEntries(slice)
    return sortedEntries(entries)
}

func populateEntries(slice []string) Entries {
    indent, indentSize := computeIndent(slice)
    fmt.Printf("[%s] %d = %d\n", indent, len(indent), indentSize)
    entries := make(Entries, 0)
    for _, item := range slice {
        i, level := 0, 0
        for strings.HasPrefix(item[i:], indent) {
            i += indentSize
            level++
        }
        key := strings.ToLower(strings.TrimSpace(item))
        addEntry(level, key, item, &entries)
    }
    return entries
}

func computeIndent(slice []string) (string, int) {
    for _, item := range slice {
        if len(item) > 0 && (item[0] == ' ' || item[0] == '\t') {
            whitespace := rune(item[0])
            for i, char := range item[1:] {
                if char != whitespace {
                    i++
                    return strings.Repeat(string(whitespace), i), i
                }
            }
        }
    }
    return "", 0
}

func addEntry(level int, key, value string, entries *Entries) {
    if level == 0 {
        *entries = append(*entries, Entry{key, value, make(Entries, 0)})
    } else {
        /*
           theEntries := *entries
           lastEntry := &theEntries[theEntries.Len()-1]
           addEntry(level-1, key, value, &lastEntry.children)
        */
        addEntry(level-1, key, value,
            &((*entries)[entries.Len()-1].children))
    }
}

func sortedEntries(entries Entries) []string {
    var indentedSlice []string
    sort.Sort(entries)
    for _, entry := range entries {
        populateIndentedStrings(entry, &indentedSlice)
    }
    return indentedSlice
}

func populateIndentedStrings(entry Entry, indentedSlice *[]string) {
    *indentedSlice = append(*indentedSlice, entry.value)
    sort.Sort(entry.children)
    for _, child := range entry.children {
        populateIndentedStrings(child, indentedSlice)
    }
}

type Entry struct {
    key      string
    value    string
    children Entries
}
type Entries []Entry

func (entries Entries) Len() int { return len(entries) }

func (entries Entries) Less(i, j int) bool {
    return entries[i].key < entries[j].key
}
func (entries Entries) Swap(i, j int) {
    entries[i], entries[j] = entries[j], entries[i]
}
