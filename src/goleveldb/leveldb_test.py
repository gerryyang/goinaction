#!/usr/bin/env python
#-*-coding: utf-8-*-

import leveldb
import os, sys

def initialize(db_name):
    db = leveldb.LevelDB(db_name);
    return db;

def insert(db, sid, name):
    db.Put(str(sid), name);

def delete(db, sid):
    db.Delete(str(sid));

def update(db, sid, name):
    db.Put(str(sid), name);

def search(db, sid):
    name = db.Get(str(sid));
    return name;

def display(db):
    info = []
    for key, value in db.RangeIter():
        print (key, value);
        info.append("%s=%s"%(key, value))
    #for k in db.RangeIter(include_value = False):
    #    print db.Get(k)
    return '\n'.join(info)

        
if __name__ == '__main__':
    db = initialize('./database/peer/org1peer0/ledgersData/stateLeveldb');
    all = display(db);
     
    print "------"
    name = search(db, 'mychannel2\x00mycc\x00a')
    print name
    update(db, 'mychannel2\x00mycc\x00a', '\x01\x02\x0030')

    print "==========="
    name = search(db, 'mychannel2\x00mycc\x00a')
    print name
