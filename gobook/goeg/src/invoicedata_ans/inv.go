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
    "encoding/binary"
    "errors"
    "fmt"
    "io"
    "strconv"
    "time"
)

type InvMarshaler struct{}

const invDateFormat = "20060102"

var byteOrder = binary.LittleEndian

func (InvMarshaler) MarshalInvoices(writer io.Writer,
    invoices []*Invoice) error {
    var write invWriterFunc = func(x interface{}) error {
        return binary.Write(writer, byteOrder, x)
    }
    if err := write(uint32(magicNumber)); err != nil {
        return err
    }
    if err := write(uint16(fileVersion)); err != nil {
        return err
    }
    if err := write(int32(len(invoices))); err != nil {
        return err
    }
    for _, invoice := range invoices {
        if err := write.writeInvoice(invoice); err != nil {
            return err
        }
    }
    return nil
}

type invWriterFunc func(interface{}) error

func (write invWriterFunc) writeInvoice(invoice *Invoice) error {
    for _, i := range []int{invoice.Id, invoice.CustomerId} {
        if err := write(int32(i)); err != nil {
            return err
        }
    }
    if err := write.writeString(invoice.DepartmentId); err != nil {
        return err
    }
    for _, date := range []time.Time{invoice.Raised, invoice.Due} {
        if err := write.writeDate(date); err != nil {
            return err
        }
    }
    if err := write.writeBool(invoice.Paid); err != nil {
        return err
    }
    if err := write.writeString(invoice.Note); err != nil {
        return err
    }
    if err := write(int32(len(invoice.Items))); err != nil {
        return err
    }
    for _, item := range invoice.Items {
        if err := write.writeItem(item); err != nil {
            return err
        }
    }
    return nil
}

func (write invWriterFunc) writeDate(date time.Time) error {
    i, err := strconv.Atoi(date.Format(invDateFormat))
    if err != nil {
        return err
    }
    return write(int32(i))
}

func (write invWriterFunc) writeBool(b bool) error {
    var v int8
    if b {
        v = 1
    }
    return write(v)
}

func (write invWriterFunc) writeString(s string) error {
    if err := write(int32(len(s))); err != nil {
        return err
    }
    return write([]byte(s))
}

func (write invWriterFunc) writeItem(item *Item) error {
    if err := write.writeString(item.Id); err != nil {
        return err
    }
    if err := write(item.Price); err != nil {
        return err
    }
    for _, i := range []int{item.Quantity, item.TaxBand} {
        if err := write(int16(i)); err != nil {
            return err
        }
    }
    return write.writeString(item.Note)
}

func (InvMarshaler) UnmarshalInvoices(reader io.Reader) ([]*Invoice,
    error) {
    version, err := checkInvVersion(reader)
    if err != nil {
        return nil, err
    }
    count, err := readIntFromInt32(reader)
    if err != nil {
        return nil, err
    }
    invoices := make([]*Invoice, 0, count)
    for i := 0; i < count; i++ {
        invoice, err := readInvInvoice(version, reader)
        if err != nil {
            return nil, err
        }
        invoices = append(invoices, invoice)
    }
    return invoices, nil
}

func readIntFromInt32(reader io.Reader) (int, error) {
    var i32 int32
    err := binary.Read(reader, byteOrder, &i32)
    return int(i32), err
}

func readIntFromInt16(reader io.Reader) (int, error) {
    var i16 int16
    err := binary.Read(reader, byteOrder, &i16)
    return int(i16), err
}

func readBoolFromInt8(reader io.Reader) (bool, error) {
    var i8 int8
    err := binary.Read(reader, byteOrder, &i8)
    return i8 == 1, err
}

func checkInvVersion(reader io.Reader) (int, error) {
    var magic uint32
    if err := binary.Read(reader, byteOrder, &magic); err != nil {
        return 0, err
    }
    if magic != magicNumber {
        return 0, errors.New("cannot read non-invoices inv file")
    }
    var version uint16
    if err := binary.Read(reader, byteOrder, &version); err != nil {
        return 0, err
    }
    if version > fileVersion {
        return 0, fmt.Errorf("version %d is too new to read", version)
    }
    return int(version), nil
}

func readInvInvoice(version int, reader io.Reader) (invoice *Invoice,
    err error) {
    invoice = &Invoice{}
    for _, i := range []*int{&invoice.Id, &invoice.CustomerId} {
        if *i, err = readIntFromInt32(reader); err != nil {
            return nil, err
        }
    }
    if version == fileVersion {
        if invoice.DepartmentId, err = readInvString(reader); err != nil {
            return nil, err
        }
    }
    for _, date := range []*time.Time{&invoice.Raised, &invoice.Due} {
        if *date, err = readInvDate(reader); err != nil {
            return nil, err
        }
    }
    if invoice.Paid, err = readBoolFromInt8(reader); err != nil {
        return nil, err
    }
    if invoice.Note, err = readInvString(reader); err != nil {
        return nil, err
    }
    var count int
    if count, err = readIntFromInt32(reader); err != nil {
        return nil, err
    }
    if invoice.Items, err = readInvItems(version, reader, count);
        err != nil {
        return nil, err
    }
    if version < fileVersion {
        updateInvoice(invoice)
    }
    return invoice, nil
}

func readInvItems(version int, reader io.Reader, count int) ([]*Item,
    error) {
    items := make([]*Item, 0, count)
    for i := 0; i < count; i++ {
        item, err := readInvItem(version, reader)
        if err != nil {
            return nil, err
        }
        items = append(items, item)
    }
    return items, nil
}

func readInvDate(reader io.Reader) (time.Time, error) {
    var n int32
    if err := binary.Read(reader, byteOrder, &n); err != nil {
        return time.Time{}, err
    }
    return time.Parse(invDateFormat, fmt.Sprint(n))
}

func readInvString(reader io.Reader) (string, error) {
    var length int32
    if err := binary.Read(reader, byteOrder, &length); err != nil {
        return "", err
    }
    raw := make([]byte, length)
    if err := binary.Read(reader, byteOrder, &raw); err != nil {
        return "", err
    }
    return string(raw), nil
}

func readInvItem(version int, reader io.Reader) (item *Item, err error) {
    item = &Item{}
    if item.Id, err = readInvString(reader); err != nil {
        return nil, err
    }
    if err = binary.Read(reader, byteOrder, &item.Price); err != nil {
        return nil, err
    }
    if item.Quantity, err = readIntFromInt16(reader); err != nil {
        return nil, err
    }
    if version == fileVersion {
        if item.TaxBand, err = readIntFromInt16(reader); err != nil {
            return nil, err
        }
    }
    if item.Note, err = readInvString(reader); err != nil {
        return nil, err
    }
    if version < fileVersion {
        if err = updateItem(item); err != nil {
            return nil, err
        }
    }
    return item, nil
}
