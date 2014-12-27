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
    "linkcheck/linkutil"
    "log"
    "net/http"
    "net/url"
    "os"
    "path/filepath"
    "regexp"
    "strings"
)

var (
    externalLinkRx *regexp.Regexp
    addChannel     chan string
    queryChannel   chan string
    seenChannel    chan bool
)

func init() {
    externalLinkRx = regexp.MustCompile("^(http|ftp|mailto):")
    // These *must* be unbuffered so that they block and properly serialize
    // access to the map
    addChannel = make(chan string)
    queryChannel = make(chan string)
    seenChannel = make(chan bool)
}

func main() {
    log.SetFlags(0)
    if len(os.Args) != 2 || os.Args[1] == "-h" || os.Args[1] == "--help" {
        log.Fatalf("usage: %s url\n", filepath.Base(os.Args[0]))
    }
    href := os.Args[1]
    if !strings.HasPrefix(href, "http://") {
        href = "http://" + href
    }
    url, err := url.Parse(href)
    if err != nil {
        log.Fatalln("- failed to read url:", err)
    }
    prepareMap()
    checkPage(href, "http://"+url.Host)
}

func prepareMap() {
    go func() {
        seen := make(map[string]bool)
        for {
            select {
            case url := <-addChannel:
                seen[url] = true
            case url := <-queryChannel:
                _, found := seen[url]
                seenChannel <- found
            }
        }
    }()
}

func alreadySeen(url string) bool {
    queryChannel <- url
    if <-seenChannel {
        return true
    }
    addChannel <- url
    return false
}

func checkPage(url, site string) {
    if alreadySeen(url) {
        return
    }
    links, err := linkutil.LinksFromURL(url)
    if err != nil {
        log.Println("-", err)
        return
    }
    fmt.Println("+ read", url)
    done := make(chan bool, len(links))
    defer close(done)
    pending := 0
    var messages []string
    for _, link := range links {
        pending += processLink(link, site, url, &messages, done)
    }
    if len(messages) > 0 {
        fmt.Println("+ links on", url)
        for _, message := range messages {
            fmt.Println("  ", message)
        }
    }
    for i := 0; i < pending; i++ {
        <-done
    }
}

func processLink(link, site, url string, messages *[]string,
    done chan<- bool) int {
    localAndParsable, link := classifyLink(link, site)
    if localAndParsable {
        go func() {
            checkPage(link, site)
            done <- true
        }()
        return 1
    }
    if message := checkExists(link, url); message != "" {
        *messages = append(*messages, message)
    }
    return 0
}

func classifyLink(link, site string) (bool, string) {
    var local, parsable bool
    url := link
    lowerLink := strings.ToLower(link)
    if strings.HasSuffix(lowerLink, ".htm") ||
        strings.HasSuffix(lowerLink, ".html") {
        parsable = true
    }
    if !externalLinkRx.MatchString(link) {
        local = true
        url = site
        if link[0] != '/' && url[len(url)-1] != '/' {
            url += "/"
        }
        url += link
    }
    return local && parsable, url
}

func checkExists(link, url string) string {
    if alreadySeen(link) {
        return ""
    }
    if _, err := http.Head(link); err != nil {
        errStr := err.Error()
        if strings.Contains(errStr, "unsupported protocol scheme") {
            return fmt.Sprint("- can't check nonhttp link: ", link)
        } else if strings.Contains(errStr, "connection timed out") {
            return fmt.Sprint("- timed out on: ", link)
        } else {
            return fmt.Sprintf("- bad link %s on page %s: %v", link, url,
                err)
        }
    }
    return fmt.Sprint("+ checked ", link)
}
