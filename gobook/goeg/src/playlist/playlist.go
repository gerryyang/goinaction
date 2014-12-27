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
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "path/filepath"
    "strconv"
    "strings"
)

type Song struct {
    Title    string
    Filename string
    Seconds  int
}

func main() {
    if len(os.Args) == 1 ||
        (!strings.HasSuffix(os.Args[1], ".pls") &&
            !strings.HasSuffix(os.Args[1], ".m3u")) {
        fmt.Printf("usage: %s <file.[pls|m3u]>\n",
            filepath.Base(os.Args[0]))
        os.Exit(1)
    }

    if rawBytes, err := ioutil.ReadFile(os.Args[1]); err != nil {
        log.Fatal(err)
    } else {
        data := string(rawBytes)
        if strings.HasSuffix(os.Args[1], ".pls") {
            songs := readPlsPlaylist(data)
            writeM3uPlaylist(songs)
        } else {
            songs := readM3uPlaylist(data)
            writePlsPlaylist(songs)
        }
    }
}

func readPlsPlaylist(data string) (songs []Song) {
    var song Song
    for _, line := range strings.Split(data, "\n") {
        if line = strings.TrimSpace(line); line == "" {
            continue
        }
        switch name, value := parsePlsLine(line); name {
        case "File":
            song.Filename = strings.Map(
                mapPlatformDirSeparator, value)
        case "Title":
            song.Title = value
        case "Length":
            var err error
            if song.Seconds, err = strconv.Atoi(value); err != nil {
                log.Printf("failed to read the duration for '%s': %v\n",
                    song.Title, err)
                song.Seconds = -1
            }
        }
        if song.Filename != "" && song.Title != "" && song.Seconds != 0 {
            songs = append(songs, song)
            song = Song{}
        }
    }
    return songs
}

func parsePlsLine(line string) (name, value string) {
    const separator = "="
    if i := strings.Index(line, separator); i > -1 {
        if j := strings.IndexAny(line, "0123456789"); j > -1 && j < i {
            name = line[:j]
            value = line[i+len(separator):]
        }
    }
    return name, value
}

func mapPlatformDirSeparator(char rune) rune {
    if char == '/' || char == '\\' {
        return filepath.Separator
    }
    return char
}

func writePlsPlaylist(songs []Song) {
    fmt.Println("[playlist]")
    for i, song := range songs {
        i++
        fmt.Printf("File%d=%s\n", i, song.Filename)
        fmt.Printf("Title%d=%s\n", i, song.Title)
        fmt.Printf("Length%d=%d\n", i, song.Seconds)
    }
    fmt.Printf("NumberOfEntries=%d\nVersion=2\n", len(songs))
}

func readM3uPlaylist(data string) (songs []Song) {
    var song Song
    for _, line := range strings.Split(data, "\n") {
        line = strings.TrimSpace(line)
        if line == "" || strings.HasPrefix(line, "#EXTM3U") {
            continue
        }
        if strings.HasPrefix(line, "#EXTINF:") {
            song.Title, song.Seconds = parseExtinfLine(line)
        } else {
            song.Filename = strings.Map(mapPlatformDirSeparator, line)
        }
        if song.Filename != "" && song.Title != "" && song.Seconds != 0 {
            songs = append(songs, song)
            song = Song{}
        }
    }
    return songs
}

func parseExtinfLine(line string) (title string, seconds int) {
    if i := strings.IndexAny(line, "-0123456789"); i > -1 {
        const separator = ","
        line = line[i:]
        if j := strings.Index(line, separator); j > -1 {
            title = line[j+len(separator):]
            var err error
            if seconds, err = strconv.Atoi(line[:j]); err != nil {
                log.Printf("failed to read the duration for '%s': %v\n",
                    title, err)
                seconds = -1
            }
        }
    }
    return title, seconds
}

func writeM3uPlaylist(songs []Song) {
    fmt.Println("#EXTM3U")
    for _, song := range songs {
        fmt.Printf("#EXTINF:%d,%s\n", song.Seconds, song.Title)
        fmt.Println(song.Filename)
    }
}
