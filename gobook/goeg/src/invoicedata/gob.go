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
//    "bytes"
    "encoding/gob"
    "errors"
    "fmt"
    "io"
//    "time"
)
/*
    // Here is how to make a custom type satisfy the gob.Encoder and
    // gob.Decoder interfaces.

type GobInvoice struct {
    Id         int
    CustomerId int
    Raised     int64 // Seconds since the Unix epoch
    Due        int64 // Seconds since the Unix epoch
    Paid       bool
    Note       string
    Items      []*Item
}

func (invoice *Invoice) GobEncode() ([]byte, error) {
    gobInvoice := GobInvoice{invoice.Id, invoice.CustomerId,
        invoice.Raised.Unix(), invoice.Due.Unix(), invoice.Paid,
        invoice.Note, invoice.Items}
    var buffer bytes.Buffer
    encoder := gob.NewEncoder(&buffer)
    err := encoder.Encode(gobInvoice)
    return buffer.Bytes(), err
}

func (invoice *Invoice) GobDecode(data []byte) error {
    var gobInvoice GobInvoice
    buffer := bytes.NewBuffer(data)
    decoder := gob.NewDecoder(buffer)
    if err := decoder.Decode(&gobInvoice); err != nil {
        return err
    }
    raised := time.Unix(gobInvoice.Raised, 0)
    due := time.Unix(gobInvoice.Due, 0)
    *invoice = Invoice{gobInvoice.Id, gobInvoice.CustomerId, raised, due,
        gobInvoice.Paid, gobInvoice.Note, gobInvoice.Items}
    return nil
}
*/
type GobMarshaler struct{}

func (GobMarshaler) MarshalInvoices(writer io.Writer,
    invoices []*Invoice) error {
    encoder := gob.NewEncoder(writer)
    if err := encoder.Encode(magicNumber); err != nil {
        return err
    }
    if err := encoder.Encode(fileVersion); err != nil {
        return err
    }
    return encoder.Encode(invoices)
}

func (GobMarshaler) UnmarshalInvoices(reader io.Reader) ([]*Invoice,
    error) {
    decoder := gob.NewDecoder(reader)
    var magic int
    if err := decoder.Decode(&magic); err != nil {
        return nil, err
    }
    if magic != magicNumber {
        return nil, errors.New("cannot read non-invoices gob file")
    }
    var version int
    if err := decoder.Decode(&version); err != nil {
        return nil, err
    }
    if version > fileVersion {
        return nil, fmt.Errorf("version %d is too new to read", version)
    }
    var invoices []*Invoice
    err := decoder.Decode(&invoices)
    return invoices, err
}
