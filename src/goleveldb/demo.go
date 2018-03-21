package main

import (
    "fmt"
    "github.com/syndtr/goleveldb/leveldb"
)

func main() {
    db, err := leveldb.OpenFile("/root/db0/ledgersData/stateLeveldb", nil)
    if err != nil {
        panic(err)
    }
    defer db.Close()

    fmt.Println("iterator")
    iter := db.NewIterator(nil, nil)
    for iter.Next() {
        fmt.Printf("key[%#v][%s] value[%#v][%s]\n", iter.Key(), iter.Key(), iter.Value(), iter.Value())
    }
        iter.Release()
        fmt.Println("----------------------------")

    data_a, _ := db.Get([]byte("mychannel\x00mycc\x00a"), nil)
    fmt.Printf("mychannelmycca[%s]\n", data_a)
    data_b, _ := db.Get([]byte("mychannel\x00mycc\x00b"), nil)
    fmt.Printf("mychannelmyccb[%s]\n", data_b)

        err = db.Put([]byte("mychannel\x00mycc\x00a"), []byte("\x01\x06\x001000"), nil)
        data_a, _ = db.Get([]byte("mychannel\x00mycc\x00a"), nil)
    fmt.Printf("new mychannelmycca[%s]\n", data_a)

        err = db.Put([]byte("mychannel\x00mycc\x00b"), []byte("\x01\x06\x002000"), nil)
        data_b, _ = db.Get([]byte("mychannel\x00mycc\x00b"), nil)
    fmt.Printf("new mychannelmyccb[%s]\n", data_b)

}