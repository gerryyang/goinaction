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
    "compress/gzip"
    "errors"
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "time"
)

const (
    magicNumber = 0x125D
    fileVersion = 101
    fileType    = "INVOICES"
    dateFormat  = "2006-01-02" // This date must always be used (see text).
)

type Invoice struct { // fileVersion
    Id           int       // 100
    CustomerId   int       // 100
    DepartmentId string    // 101
    Raised       time.Time // 100
    Due          time.Time // 100
    Paid         bool      // 100
    Note         string    // 100
    Items        []*Item   // 100
}

type Item struct { // fileVersion
    Id       string  // 100
    Price    float64 // 100
    Quantity int     // 100
    TaxBand  int     // 101
    Note     string  // 100
}

type InvoicesMarshaler interface {
    MarshalInvoices(writer io.Writer, invoices []*Invoice) error
}

type InvoicesUnmarshaler interface {
    UnmarshalInvoices(reader io.Reader) ([]*Invoice, error)
}

func main() {
    log.SetFlags(0)
    if len(os.Args) != 3 || os.Args[1] == "-h" || os.Args[1] == "--help" {
        log.Fatalf("usage: %s infile.ext outfile.ext\n"+
            ".ext may be any of .gob, .inv, .jsn, .json, .txt, "+
            "or .xml, optionally gzipped (e.g., .gob.gz)\n",
            filepath.Base(os.Args[0]))
    }
    inFilename, outFilename := os.Args[1], os.Args[2]
    if inFilename == outFilename {
        log.Fatalln("won't overwrite a file with itself")
    }

    invoices, err := readInvoiceFile(inFilename)
    if err != nil {
        log.Fatalln("Failed to read:", err)
    }
    if err := writeInvoiceFile(outFilename, invoices); err != nil {
        log.Fatalln("Failed to write:", err)
    }
}

func readInvoiceFile(filename string) ([]*Invoice, error) {
    file, closer, err := openInvoiceFile(filename)
    if closer != nil {
        defer closer()
    }
    if err != nil {
        return nil, err
    }
    return readInvoices(file, suffixOf(filename))
}

func openInvoiceFile(filename string) (io.ReadCloser, func(), error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, nil, err
    }
    closer := func() { file.Close() }
    var reader io.ReadCloser = file
    var decompressor *gzip.Reader
    if strings.HasSuffix(filename, ".gz") {
        if decompressor, err = gzip.NewReader(file); err != nil {
            return file, closer, err
        }
        closer = func() { decompressor.Close(); file.Close() }
        reader = decompressor
    }
    return reader, closer, nil
}

func readInvoices(reader io.Reader, suffix string) ([]*Invoice, error) {
    var unmarshaler InvoicesUnmarshaler
    switch suffix {
    case ".gob":
        unmarshaler = GobMarshaler{}
    case ".inv":
        unmarshaler = InvMarshaler{}
    case ".jsn", ".json":
        unmarshaler = JSONMarshaler{}
    case ".txt":
        unmarshaler = TxtMarshaler{}
    case ".xml":
        unmarshaler = XMLMarshaler{}
    }
    if unmarshaler != nil {
        return unmarshaler.UnmarshalInvoices(reader)
    }
    return nil, fmt.Errorf("unrecognized input suffix: %s", suffix)
}

func writeInvoiceFile(filename string, invoices []*Invoice) error {
    file, closer, err := createInvoiceFile(filename)
    if closer != nil {
        defer closer()
    }
    if err != nil {
        return err
    }
    return writeInvoices(file, suffixOf(filename), invoices)
}

func createInvoiceFile(filename string) (io.WriteCloser, func(), error) {
    file, err := os.Create(filename)
    if err != nil {
        return nil, nil, err
    }
    closer := func() { file.Close() }
    var writer io.WriteCloser = file
    var compressor *gzip.Writer
    if strings.HasSuffix(filename, ".gz") {
        compressor = gzip.NewWriter(file)
        closer = func() { compressor.Close(); file.Close() }
        writer = compressor
    }
    return writer, closer, nil
}

func writeInvoices(writer io.Writer, suffix string,
    invoices []*Invoice) error {
    var marshaler InvoicesMarshaler
    switch suffix {
    case ".gob":
        marshaler = GobMarshaler{}
    case ".inv":
        marshaler = InvMarshaler{}
    case ".jsn", ".json":
        marshaler = JSONMarshaler{}
    case ".txt":
        marshaler = TxtMarshaler{}
    case ".xml":
        marshaler = XMLMarshaler{}
    }
    if marshaler != nil {
        return marshaler.MarshalInvoices(writer, invoices)
    }
    return errors.New("unrecognized output suffix")
}

func update(invoices []*Invoice) error {
    for _, invoice := range invoices {
        updateInvoice(invoice)
        for _, item := range invoice.Items {
            if err := updateItem(item); err != nil {
                return err
            }
        }
    }
    return nil
}

func suffixOf(filename string) string {
    suffix := filepath.Ext(filename)
    if suffix == ".gz" {
        suffix = filepath.Ext(filename[:len(filename)-3])
    }
    return suffix
}

func updateInvoice(invoice *Invoice) {
    switch {
    case invoice.Id < 3000:
        invoice.DepartmentId = "GEN"
    case invoice.Id < 4000:
        invoice.DepartmentId = "MKT"
    case invoice.Id < 5000:
        invoice.DepartmentId = "COM"
    case invoice.Id < 6000:
        invoice.DepartmentId = "EXP"
    case invoice.Id < 7000:
        invoice.DepartmentId = "INP"
    case invoice.Id < 8000:
        invoice.DepartmentId = "TZZ"
    case invoice.Id < 9000:
        invoice.DepartmentId = "V20"
    default:
        invoice.DepartmentId = "X15"
    }
}

func updateItem(item *Item) (err error) {
    if item.TaxBand, err = strconv.Atoi(item.Id[2:3]); err != nil {
        return fmt.Errorf("invalid item ID: %s", item.Id)
    }
    return nil
}
