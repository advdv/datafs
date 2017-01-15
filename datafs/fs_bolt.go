package datafs

import (
	"github.com/boltdb/bolt"
)

//BoltFS creates a file system on top of the bolt memory-map kv database
type BoltFS struct {
	db *bolt.DB

	*EmptyFS //@TODO progressively make remove this
}

//NewBoltFS will setup the database for the fs
func NewBoltFS(db *bolt.DB) *BoltFS {
	return &BoltFS{
		db:      db,
		EmptyFS: &EmptyFS{},
	}
}
