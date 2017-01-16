package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/advanderveer/datafs/datafs"
	"github.com/boltdb/bolt"
	"github.com/keybase/kbfs/dokan"
)

func main() {
	log.Printf("started")
	defer log.Printf("exited")
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	tmpdir, err := ioutil.TempDir("", "datafs_")
	if err != nil {
		log.Fatal(err)
	}

	db, err := bolt.Open(filepath.Join(tmpdir, "fs.bolt"), 0777, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("using bolt db '%s' as filesystem backend", db.Path())
	defer db.Close()

	fs, err := datafs.NewBoltFS(db)
	if err != nil {
		log.Fatal(err)
	}

	conf := &dokan.Config{
		FileSystem: fs,
		Path:       `T:\`,
	}

	mnt, err := dokan.Mount(conf)
	if err != nil {
		log.Fatal(err)
	}

	defer mnt.Close()
	<-sigCh
}
