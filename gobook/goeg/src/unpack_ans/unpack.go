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
    "archive/tar"
    "archive/zip"
    "compress/bzip2"
    "compress/gzip"
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "strings"
)

func main() {
    log.SetFlags(0)
    if len(os.Args) == 1 || os.Args[1] == "-h" || os.Args[1] == "--help" {
        fmt.Printf("usage: %s archive.{zip,tar,tar.gz,tar.bz2}\n",
            filepath.Base(os.Args[0]))
        os.Exit(1)

    }
    filename := os.Args[1]
    if !validSuffix(filename) {
        log.Fatalln("unrecognized archive suffix")
    }
    if err := unpackArchive(filename); err != nil {
        log.Fatalln(err)
    }
}

func validSuffix(filename string) bool {
    for _, suffix := range []string{".zip", ".tar", ".tar.gz",
        ".tar.bz2"} {
        if strings.HasSuffix(filename, suffix) {
            return true
        }
    }
    return false
}

func unpackArchive(filename string) error {
    if strings.HasSuffix(filename, ".zip") {
        return unpackZip(filename)
    }
    return unpackTar(filename)
}

func unpackZip(filename string) (err error) {
    var reader *zip.ReadCloser
    if reader, err = zip.OpenReader(filename); err != nil {
        return err
    }
    defer reader.Close()
    for _, zipFile := range reader.Reader.File {
        filename := sanitizedName(zipFile.Name)
        if strings.HasSuffix(zipFile.Name, "/") ||
            strings.HasSuffix(zipFile.Name, "\\") {
            if err = os.MkdirAll(filename, 0755); err != nil {
                return err
            }
        } else {
            if err = unpackZippedFile(filename, zipFile); err != nil {
                return err
            }
        }
    }
    return nil
}

func unpackZippedFile(filename string, zipFile *zip.File) (err error) {
    var writer *os.File
    if writer, err = os.Create(filename); err != nil {
        return err
    }
    defer writer.Close()
    var reader io.ReadCloser
    if reader, err = zipFile.Open(); err != nil {
        return err
    }
    defer reader.Close()
    if _, err = io.Copy(writer, reader); err != nil {
        return err
    }
    if filename == zipFile.Name {
        fmt.Println(filename)
    } else {
        fmt.Printf("%s [%s]\n", filename, zipFile.Name)
    }
    return nil
}

func unpackTar(filename string) (err error) {
    var file *os.File
    if file, err = os.Open(filename); err != nil {
        return err
    }
    defer file.Close()
    var fileReader io.Reader = file
    var decompressor *gzip.Reader
    if strings.HasSuffix(filename, ".gz") {
        if decompressor, err = gzip.NewReader(file); err != nil {
            return err
        }
        defer decompressor.Close()
    } else if strings.HasSuffix(filename, ".bz2") {
        fileReader = bzip2.NewReader(file)
    }
    var reader *tar.Reader
    if decompressor != nil {
        reader = tar.NewReader(decompressor)
    } else {
        reader = tar.NewReader(fileReader)
    }
    return unpackTarFiles(reader)
}

func unpackTarFiles(reader *tar.Reader) (err error) {
    var header *tar.Header
    for {
        if header, err = reader.Next(); err != nil {
            if err == io.EOF {
                return nil // OK
            }
            return err
        }
        filename := sanitizedName(header.Name)
        switch header.Typeflag {
        case tar.TypeDir:
            if err = os.MkdirAll(filename, 0755); err != nil {
                return err
            }
        case tar.TypeReg:
            if err = unpackTarFile(filename, header.Name, reader); err != nil {
                return err
            }
        }
    }
    return nil
}

func unpackTarFile(filename, tarFilename string,
    reader *tar.Reader) (err error) {
    var writer *os.File
    if writer, err = os.Create(filename); err != nil {
        return err
    }
    defer writer.Close()
    if _, err = io.Copy(writer, reader); err != nil {
        return err
    }
    if filename == tarFilename {
        fmt.Println(filename)
    } else {
        fmt.Printf("%s [%s]\n", filename, tarFilename)
    }
    return nil
}

func sanitizedName(filename string) string {
    if len(filename) > 1 && filename[1] == ':' {
        filename = filename[2:]
    }
    filename = strings.TrimLeft(filename, "\\/.")
    filename = strings.Replace(filename, "../", "", -1)
    return strings.Replace(filename, "..\\", "", -1)
}
