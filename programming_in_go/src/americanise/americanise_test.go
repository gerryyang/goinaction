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
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
)

const (
	inFilename       = "input.txt"
	expectedFilename = "expected.txt"
	actualFilename   = "actual.txt"
)

func TestAmericanize(t *testing.T) {
	log.SetFlags(0)
	log.Println("TEST americanize")

	path, _ := filepath.Split(os.Args[0])
	var inFile, outFile *os.File
	var err error

	inFilename := filepath.Join(path, inFilename)
	if inFile, err = os.Open(inFilename); err != nil {
		t.Fatal(err)
	}
	defer inFile.Close()

	outFilename := filepath.Join(path, actualFilename)
	if outFile, err = os.Create(outFilename); err != nil {
		t.Fatal(err)
	}
	defer outFile.Close()
	defer os.Remove(outFilename)

	if err := americanise(inFile, outFile); err != nil {
		t.Fatal(err)
	}

	compare(outFilename, filepath.Join(path, expectedFilename), t)
}

func compare(actual, expected string, t *testing.T) {

	if actualBytes, err := ioutil.ReadFile(actual); err != nil {
		t.Fatal(err)
	} else if expectedBytes, err := ioutil.ReadFile(expected); err != nil {
		t.Fatal(err)
	} else {
		if bytes.Compare(actualBytes, expectedBytes) != 0 {
			t.Fatal("actual != expected")
		}
	}
}
