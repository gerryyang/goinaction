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
    "encoding/xml"
    "fmt"
    "io"
    "strings"
    "time"
)

type XMLMarshaler struct{}

type XMLInvoices struct {
    XMLName xml.Name      `xml:"INVOICES"`
    Version int           `xml:"version,attr"`
    Invoice []*XMLInvoice `xml:"INVOICE"`
}

type XMLInvoice struct {
    XMLName    xml.Name   `xml:"INVOICE"`
    Id         int        `xml:",attr"`
    CustomerId int        `xml:",attr"`
    Raised     string     `xml:",attr"`
    Due        string     `xml:",attr"`
    Paid       bool       `xml:",attr"`
    Note       string     `xml:"NOTE"`
    Item       []*XMLItem `xml:"ITEM"`
}

type XMLItem struct {
    XMLName  xml.Name `xml:"ITEM"`
    Id       string   `xml:",attr"`
    Price    float64  `xml:",attr"`
    Quantity int      `xml:",attr"`
    Note     string   `xml:"NOTE"`
}


func XMLInvoicesForInvoices(invoices []*Invoice) *XMLInvoices {
    xmlInvoices := &XMLInvoices{
        Version: fileVersion,
        Invoice: make([]*XMLInvoice, 0, len(invoices)),
    }
    for _, invoice := range invoices {
        xmlInvoices.Invoice = append(xmlInvoices.Invoice,
            XMLInvoiceForInvoice(invoice))
    }
    return xmlInvoices
}


func XMLInvoiceForInvoice(invoice *Invoice) *XMLInvoice {
    xmlInvoice := &XMLInvoice{
        Id:         invoice.Id,
        CustomerId: invoice.CustomerId,
        Raised:     invoice.Raised.Format(dateFormat),
        Due:        invoice.Due.Format(dateFormat),
        Paid:       invoice.Paid,
        Note:       invoice.Note,
        Item:       make([]*XMLItem, 0, len(invoice.Items)),
    }
    for _, item := range invoice.Items {
        xmlItem := &XMLItem{
            Id:       item.Id,
            Price:    item.Price,
            Quantity: item.Quantity,
            Note:     item.Note,
        }
        xmlInvoice.Item = append(xmlInvoice.Item, xmlItem)
    }
    return xmlInvoice
}


func (xmlInvoices *XMLInvoices) Invoices() (invoices []*Invoice,
    err error) {
    invoices = make([]*Invoice, 0, len(xmlInvoices.Invoice))
    for _, xmlInvoice := range xmlInvoices.Invoice {
        invoice, err := xmlInvoice.Invoice()
        if err != nil {
            return nil, err
        }
        invoices = append(invoices, invoice)
    }
    return invoices, nil
}


func (xmlInvoice *XMLInvoice) Invoice() (invoice *Invoice, err error) {
    invoice = &Invoice{
        Id:         xmlInvoice.Id,
        CustomerId: xmlInvoice.CustomerId,
        Paid:       xmlInvoice.Paid,
        Note:       strings.TrimSpace(xmlInvoice.Note),
        Items:      make([]*Item, 0, len(xmlInvoice.Item)),
    }
    if invoice.Raised, err = time.Parse(dateFormat, xmlInvoice.Raised);
        err != nil {
        return nil, err
    }
    if invoice.Due, err = time.Parse(dateFormat, xmlInvoice.Due);
        err != nil {
        return nil, err
    }
    for _, xmlItem := range xmlInvoice.Item {
        item := &Item{
            Id:       xmlItem.Id,
            Price:    xmlItem.Price,
            Quantity: xmlItem.Quantity,
            Note:     strings.TrimSpace(xmlItem.Note),
        }
        invoice.Items = append(invoice.Items, item)
    }
    return invoice, nil
}


func (XMLMarshaler) MarshalInvoices(writer io.Writer,
    invoices []*Invoice) error {
    if _, err := writer.Write([]byte(xml.Header)); err != nil {
        return err
    }
    xmlInvoices := XMLInvoicesForInvoices(invoices)
    encoder := xml.NewEncoder(writer)
    return encoder.Encode(xmlInvoices)
}


func (XMLMarshaler) UnmarshalInvoices(reader io.Reader) ([]*Invoice,
    error) {
    xmlInvoices := &XMLInvoices{}
    decoder := xml.NewDecoder(reader)
    if err := decoder.Decode(xmlInvoices); err != nil {
        return nil, err
    }
    if xmlInvoices.Version > fileVersion {
        return nil, fmt.Errorf("version %d is too new to read",
            xmlInvoices.Version)
    }
    return xmlInvoices.Invoices()
}
