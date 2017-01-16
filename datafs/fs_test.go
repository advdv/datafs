package datafs_test

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/advanderveer/datafs/datafs"
	"github.com/boltdb/bolt"
	"github.com/keybase/kbfs/dokan"
)

var mntpath = `T:\`

type tester interface {
	Fatalf(format string, a ...interface{})
}

func path(parts ...string) string {
	p := append([]string{mntpath}, parts...)
	return filepath.Join(p...)
}

func testfs(t tester) *datafs.BoltFS {
	logs := log.New(os.Stderr, "datafs/", log.Lshortfile)
	tmpdir, err := ioutil.TempDir("", "dfs_test_")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	db, err := bolt.Open(filepath.Join(tmpdir, "fs.bolt"), 0666, nil)
	if err != nil {
		t.Fatalf("failed to open bolt db: %v", err)
	}

	fs, err := datafs.NewBoltFS(logs, db)
	if err != nil {
		t.Fatalf("failed to create fs: %v", err)
	}

	return fs
}

func TestBasicOperations(t *testing.T) {
	fs := testfs(t)
	mnt, err := dokan.Mount(&dokan.Config{FileSystem: fs, Path: mntpath})
	if err != nil {
		t.Fatal(err)
	}

	cases := map[string]func(string, *testing.T){
		"ReadNonExistingFile":   CaseOpenNonExistingFile,
		"CreateNonExistingFile": CaseCreateNonExistingFile,
		"ReadCreatedFile":       CaseReadCreatedFile,
	}
	for name, fn := range cases {
		t.Run(name, func(t *testing.T) {
			fn(mntpath, t)
		})
	}

	err = mnt.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func CaseOpenNonExistingFile(p string, t *testing.T) {
	_, err := ioutil.ReadFile(filepath.Join(p, "abc.txt"))
	if !os.IsNotExist(err) {
		t.Errorf("expected a non existing file error, got: %v", err)
	}
}

func CaseCreateNonExistingFile(p string, t *testing.T) {
	err := ioutil.WriteFile(filepath.Join(p, "abc.txt"), []byte("hello, world"), 0777)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func CaseReadCreatedFile(p string, t *testing.T) {
	input := []byte("hello, world")
	fpath := filepath.Join(p, "abc.txt")
	err := ioutil.WriteFile(fpath, input, 0777)
	if err != nil {
		t.Error(err)
	}

	output, err := ioutil.ReadFile(fpath)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(output, input) {
		t.Errorf("output(len %d) should equal input(len %d)", len(output), len(input))
	}
}
