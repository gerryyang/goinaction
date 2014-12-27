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
    "bufio"
    "errors"
    "fmt"
    "io"
    "strings"
    "time"
)

const noteSep = ":"

type TxtMarshaler struct{}

func (TxtMarshaler) MarshalInvoices(writer io.Writer,
    invoices []*Invoice) error {
    bufferedWriter := bufio.NewWriter(writer)
    defer bufferedWriter.Flush()
    var write writerFunc = func(format string,
        args ...interface{}) error {
        _, err := fmt.Fprintf(bufferedWriter, format, args...)
        return err
    }
    if err := write("%s %d\n", fileType, fileVersion); err != nil {
        return err
    }
    for _, invoice := range invoices {
        if err := write.writeInvoice(invoice); err != nil {
            return err
        }
    }
    return nil
}

type writerFunc func(string, ...interface{}) error

func (write writerFunc) writeInvoice(invoice *Invoice) error {
    note := ""
    if invoice.Note != "" {
        note = noteSep + " " + invoice.Note
    }
    if err := write("INVOICE ID=%d CUSTOMER=%d DEPARTMENT=%s RAISED=%s "+
        "DUE=%s PAID=%t%s\n", invoice.Id, invoice.CustomerId,
        invoice.DepartmentId, invoice.Raised.Format(dateFormat),
        invoice.Due.Format(dateFormat), invoice.Paid, note); err != nil {
        return err
    }
    if err := write.writeItems(invoice.Items); err != nil {
        return err
    }
    return write("\f\n")
}

func (write writerFunc) writeItems(items []*Item) error {
    for _, item := range items {
        note := ""
        if item.Note != "" {
            note = noteSep + " " + item.Note
        }
        if err := write("ITEM ID=%s PRICE=%.2f QUANTITY=%d TAXBAND=%d%s\n",
            item.Id, item.Price, item.Quantity, item.TaxBand, note);
            err != nil {
            return err
        }
    }
    return nil
}

func (TxtMarshaler) UnmarshalInvoices(reader io.Reader) (
    invoices []*Invoice, err error) {
    bufferedReader := bufio.NewReader(reader)
    var version int
    if version, err = checkTxtVersion(bufferedReader); err != nil {
        return nil, err
    }
    var line string
    eof := false
    for lino := 2; !eof; lino++ {
        line, err = bufferedReader.ReadString('\n')
        if err == io.EOF {
            err = nil
            eof = true
        } else if err != nil {
            return nil, err
        }
        if invoices, err = parseTxtLine(version, lino, line, invoices);
            err != nil {
            return nil, err
        }
    }
    return invoices, nil
}

func checkTxtVersion(bufferedReader *bufio.Reader) (version int,
    err error) {
    if _, err = fmt.Fscanf(bufferedReader, "INVOICES %d\n", &version);
        err != nil {
        err = errors.New("cannot read non-invoices text file")
    } else if version > fileVersion {
        err = fmt.Errorf("version %d is too new to read", version)
    }
    return version, err
}

func parseTxtLine(version, lino int, line string, invoices []*Invoice) (
    []*Invoice, error) {
    var err error
    if strings.HasPrefix(line, "INVOICE") {
        var invoice *Invoice
        invoice, err = parseTxtInvoice(version, lino, line)
        invoices = append(invoices, invoice)
    } else if strings.HasPrefix(line, "ITEM") {
        if len(invoices) == 0 {
            err = fmt.Errorf("item outside of an invoice line %d", lino)
        } else {
            var item *Item
            item, err = parseTxtItem(version, lino, line)
            items := &invoices[len(invoices)-1].Items
            *items = append(*items, item)
        }
    }
    return invoices, err
}

func parseTxtInvoice(version, lino int, line string) (invoice *Invoice,
    err error) {
    invoice = &Invoice{}
    var raised, due string
    if version == fileVersion {
        if _, err = fmt.Sscanf(line, "INVOICE ID=%d CUSTOMER=%d "+
            "DEPARTMENT=%s RAISED=%s DUE=%s PAID=%t", &invoice.Id,
            &invoice.CustomerId, &invoice.DepartmentId, &raised, &due,
            &invoice.Paid); err != nil {
            return nil, fmt.Errorf("invalid invoice %v line %d", err, lino)
        }
    } else {
        if _, err = fmt.Sscanf(line, "INVOICE ID=%d CUSTOMER=%d "+
            "RAISED=%s DUE=%s PAID=%t", &invoice.Id,
            &invoice.CustomerId, &raised, &due, &invoice.Paid);
                err != nil {
            return nil, fmt.Errorf("invalid invoice %v line %d", err, lino)
        }
    }
    if invoice.Raised, err = time.Parse(dateFormat, raised); err != nil {
        return nil, fmt.Errorf("invalid raised %v line %d", err, lino)
    }
    if invoice.Due, err = time.Parse(dateFormat, due); err != nil {
        return nil, fmt.Errorf("invalid due %v line %d", err, lino)
    }
    if i := strings.Index(line, noteSep); i > -1 {
        invoice.Note = strings.TrimSpace(line[i+len(noteSep):])
    }
    if version < fileVersion {
        updateInvoice(invoice)
    }
    return invoice, nil
}

func parseTxtItem(version, lino int, line string) (item *Item,
    err error) {
    item = &Item{}
    if version == fileVersion {
        if _, err = fmt.Sscanf(line, "ITEM ID=%s PRICE=%f "+
            "QUANTITY=%d TAXBAND=%d", &item.Id, &item.Price,
            &item.Quantity, &item.TaxBand); err != nil {
            return nil, fmt.Errorf("invalid item %v line %d", err, lino)
        }
    } else {
        if _, err = fmt.Sscanf(line, "ITEM ID=%s PRICE=%f QUANTITY=%d",
            &item.Id, &item.Price, &item.Quantity); err != nil {
            return nil, fmt.Errorf("invalid item %v line %d", err, lino)
        }
    }
    if i := strings.Index(line, noteSep); i > -1 {
        item.Note = strings.TrimSpace(line[i+len(noteSep):])
    }
    if version < fileVersion {
        updateItem(item)
    }
    return item, nil
}
