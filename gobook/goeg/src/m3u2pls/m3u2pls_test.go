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
    "io/ioutil"
    "log"
    "os"
    "testing"
)

func TestReadM3uPlaylist(t *testing.T) {
    log.SetFlags(0)
    log.Println("TEST m3u2pls")

    songs := readM3uPlaylist(M3U)
    for i, song := range songs {
        if song.Title != ExpectedSongs[i].Title {
            t.Fatalf("%q != %q", song.Title, ExpectedSongs[i].Title)
        }
        if song.Filename != ExpectedSongs[i].Filename {
            t.Fatalf("%q != %q", song.Filename,
                ExpectedSongs[i].Filename)
        }
        if song.Seconds != ExpectedSongs[i].Seconds {
            t.Fatalf("%d != %d", song.Seconds,
                ExpectedSongs[i].Seconds)
        }
    }
}

func TestWritePlsPlaylist(t *testing.T) {
    songs := readM3uPlaylist(M3U)
    var err error
    reader, writer := os.Stdin, os.Stdout
    if reader, writer, err = os.Pipe(); err != nil {
        t.Fatal(err)
    }
    os.Stdout = writer
    writePlsPlaylist(songs)
    writer.Close()
    actual, err := ioutil.ReadAll(reader)
    if err != nil {
        t.Fatal(err)
    }
    reader.Close()
    if string(actual) != ExpectedPls {
        t.Fatal("actual != expected")
    }
}

const M3U = `#EXTM3U
#EXTINF:315,David Bowie - Space Oddity
Music/David Bowie/Singles 1/01-Space Oddity.ogg
#EXTINF:-1,David Bowie - Changes
Music/David Bowie/Singles 1/02-Changes.ogg
#EXTINF:258,David Bowie - Starman
Music/David Bowie/Singles 1/03-Starman.ogg`

var ExpectedSongs = []Song{
    {"David Bowie - Space Oddity",
        "Music/David Bowie/Singles 1/01-Space Oddity.ogg", 315},
    {"David Bowie - Changes",
        "Music/David Bowie/Singles 1/02-Changes.ogg", -1},
    {"David Bowie - Starman",
        "Music/David Bowie/Singles 1/03-Starman.ogg", 258},
    }

var ExpectedPls = `[playlist]
File1=Music/David Bowie/Singles 1/01-Space Oddity.ogg
Title1=David Bowie - Space Oddity
Length1=315
File2=Music/David Bowie/Singles 1/02-Changes.ogg
Title2=David Bowie - Changes
Length2=-1
File3=Music/David Bowie/Singles 1/03-Starman.ogg
Title3=David Bowie - Starman
Length3=258
NumberOfEntries=3
Version=2
`
