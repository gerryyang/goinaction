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

package safeslice

type safeSlice chan commandData

type commandData struct {
    action  commandAction
    index   int
    item    interface{}
    result  chan<- interface{}
    data    chan<- []interface{}
    updater UpdateFunc
}

type commandAction int

const (
    insert commandAction = iota
    remove
    at
    update
    end
    length
)

type UpdateFunc func(interface{}) interface{}

type SafeSlice interface {
    Append(interface{})     // Append the given item to the slice
    At(int) interface{}     // Return the item at the given index position
    Close() []interface{}   // Close the channel and return the slice
    Delete(int)             // Delete the item at the given index position
    Len() int               // Return the number of items in the slice
    Update(int, UpdateFunc) // Update the item at the given index position
}

func New() SafeSlice {
    slice := make(safeSlice)
    go slice.run()
    return slice
}

func (slice safeSlice) run() {
    list := make([]interface{}, 0)
    for command := range slice {
        switch command.action {
        case insert:
            list = append(list, command.item)
        case remove: // potentially expensive for long lists
            if 0 <= command.index && command.index < len(list) {
                list = append(list[:command.index],
                    list[command.index+1:]...)
            }
        case at:
            if 0 <= command.index && command.index < len(list) {
                command.result <- list[command.index]
            } else {
                command.result <- nil
            }
        case length:
            command.result <- len(list)
        case update:
            if 0 <= command.index && command.index < len(list) {
                list[command.index] = command.updater(list[command.index])
            }
        case end:
            close(slice)
            command.data <- list
        }
    }
}

func (slice safeSlice) Append(item interface{}) {
    slice <- commandData{action: insert, item: item}
}

func (slice safeSlice) Delete(index int) {
    slice <- commandData{action: remove, index: index}
}

func (slice safeSlice) At(index int) interface{} {
    reply := make(chan interface{})
    slice <- commandData{at, index, nil, reply, nil, nil}
    return <-reply
}

func (slice safeSlice) Len() int {
    reply := make(chan interface{})
    slice <- commandData{action: length, result: reply}
    return (<-reply).(int)
}

// If the updater calls a safeSlice method we will get deadlock!
func (slice safeSlice) Update(index int, updater UpdateFunc) {
    slice <- commandData{action: update, index: index, updater: updater}
}

func (slice safeSlice) Close() []interface{} {
    reply := make(chan []interface{})
    slice <- commandData{action: end, data: reply}
    return <-reply
}
