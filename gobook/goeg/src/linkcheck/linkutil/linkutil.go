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

package linkutil

import (
    "fmt"
    "io"
    "io/ioutil"
    "net/http"
    "regexp"
)

var hrefRx *regexp.Regexp

func init() {
    hrefRx = regexp.MustCompile(`<a[^>]+href=['"]?([^'">]+)['"]?`)
}

func LinksFromURL(url string) ([]string, error) {
    response, err := http.Get(url)
    if err != nil {
        return nil, fmt.Errorf("failed to get page: %s", err)
    }
    links, err := LinksFromReader(response.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to parse page: %s", err)
    }
    return links, nil
}

func LinksFromReader(reader io.Reader) ([]string, error) {
    html, err := ioutil.ReadAll(reader)
    if err != nil {
        return nil, err
    }
    // FindAllSubmatch returns a slice of slices of slices!
    // The outer level is each match, the next level is the groups, 0 for
    // the whole match 1 for first set of ()s etc; the innermost level
    // being individual () matches
    uniqueLinks := make(map[string]bool)
    for _, submatch := range hrefRx.FindAllSubmatch(html, -1) {
        uniqueLinks[string(submatch[1])] = true
    }
    links := make([]string, len(uniqueLinks))
    i := 0
    for link := range uniqueLinks {
        links[i] = link
        i++
    }
    return links, nil
}
