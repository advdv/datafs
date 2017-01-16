package datafs

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/boltdb/bolt"
	"github.com/keybase/kbfs/dokan"
	"golang.org/x/net/context"
)

var (
	//BucketNameMetadata is the bucket name that holds filesystem metadata
	BucketNameMetadata = []byte("metadata")
)

//BoltFile is a file that is persisted in a memory mapped file instead of a block device
type BoltFile struct {
	IsDirectory bool `json:"d"`

	EmptyFile
}

//NewBoltFile sets up memory for a boltfile
func NewBoltFile(isdir bool) *BoltFile {
	return &BoltFile{
		IsDirectory: isdir,
	}
}

//LoadBoltFile will attempt to read and deserialize a file from the database
func LoadBoltFile(b *bolt.Bucket, path string) (f *BoltFile, err error) {
	data := b.Get([]byte(path))
	if data == nil {
		return nil, os.ErrNotExist
	}

	f = &BoltFile{}
	err = json.Unmarshal(data, f)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize file '%s': %v", path, err)
	}

	return f, nil
}

//Save the boltfile state to the database
func (f *BoltFile) Save(b *bolt.Bucket, path string) error {
	data, err := json.Marshal(f)
	if err != nil {
		return fmt.Errorf("failed to serialize file '%s': %v", path, err)
	}

	return b.Put([]byte(path), data)
}

//IsDir returns if the metadata information describes a file
func (f *BoltFile) IsDir() bool {
	//@TODO implement
	return false
}

// FindFiles is the readdir. The function is a callback that should be called
// with each file. The same NamedStat may be reused for subsequent calls.
//
// Pattern will be an empty string unless UseFindFilesWithPattern is enabled - then
// it may be a pattern like `*.png` to match. All implementations must be prepared
// to handle empty strings as patterns.
func (f *BoltFile) FindFiles(ctx context.Context, fi *dokan.FileInfo, pattern string, fillStatCallback func(*dokan.NamedStat) error) error {

	// &dokan.NamedStat{}
	return nil
}

// GetFileInformation - corresponds to stat.
func (f *BoltFile) GetFileInformation(ctx context.Context, fi *dokan.FileInfo) (st *dokan.Stat, err error) {
	st = &dokan.Stat{
		Creation:           time.Now(),                // Timestamps for the file
		LastAccess:         time.Now(),                // Timestamps for the file
		LastWrite:          time.Now(),                // Timestamps for the file
		FileSize:           5,                         // FileSize is the size of the file in bytes
		FileIndex:          1000,                      // FileIndex is a 64 bit (nearly) unique ID of the file
		FileAttributes:     dokan.FileAttributeNormal, // FileAttributes bitmask holds the file attributes
		VolumeSerialNumber: 0,                         // VolumeSerialNumber is the serial number of the volume (0 is fine)
		NumberOfLinks:      1,                         // NumberOfLinks can be omitted, if zero set to 1.
		ReparsePointTag:    0,                         // ReparsePointTag is for WIN32_FIND_DATA dwReserved0 for reparse point tags, typically it can be omitted.
	}

	if f.IsDir() {
		st.FileAttributes = dokan.FileAttributeDirectory
	}

	return st, nil
}

//BoltFS creates a file system on top of the bolt memory-map kv database
type BoltFS struct {
	logs *log.Logger
	db   *bolt.DB

	*EmptyFS //@TODO progressively make remove this
}

//NewBoltFS will setup the database for the fs
func NewBoltFS(logs *log.Logger, db *bolt.DB) (fs *BoltFS, err error) {
	fs = &BoltFS{
		logs:    logs,
		db:      db,
		EmptyFS: &EmptyFS{},
	}

	if err = fs.db.Update(func(tx *bolt.Tx) error {
		b, txerr := tx.CreateBucketIfNotExists(BucketNameMetadata)
		if txerr != nil {
			return txerr
		}

		root := NewBoltFile(true)
		txerr = root.Save(b, `\`)
		if txerr != nil {
			return fmt.Errorf("failed to create root: %v", txerr)
		}

		return txerr
	}); err != nil {
		return nil, err
	}

	return fs, nil
}

// GetVolumeInformation returns information about the volume.
func (fs *BoltFS) GetVolumeInformation(ctx context.Context) (dokan.VolumeInformation, error) {
	fs.logs.Printf("BoltFS.GetVolumeInformation(ctx)")
	return dokan.VolumeInformation{
		//Maximum file name component length, in bytes, supported by the specified file system. A file name component is that portion of a file name between backslashes.
		MaximumComponentLength: 0xFF, // This can be changed.
		FileSystemFlags: dokan.FileCasePreservedNames | //The file system preserves the case of file names when it places a name on disk.
			dokan.FileCaseSensitiveSearch | //The file system supports case-sensitive file names.
			dokan.FileUnicodeOnDisk | //The file system supports Unicode in file names.
			dokan.FileSupportsReparsePoints | //The file system supports reparse points.
			dokan.FileSupportsRemoteStorage, //The file system supports remote storage.
		FileSystemName: "Nerdalize Compute Engine",
		VolumeName:     "My-Organization",
	}, nil
}

// CreateFile is called to open and create files.
func (fs *BoltFS) CreateFile(ctx context.Context, fi *dokan.FileInfo, cd *dokan.CreateData) (f dokan.File, isDir bool, err error) {
	fs.logs.Printf("BoltFS.CreateFile(ctx, fi{Path: '%s'} cd{CreateDisposition: '%d'})", fi.Path(), cd.CreateDisposition)

	//Specifies what to do, depending on whether the file already exists, as one of the following values.
	switch cd.CreateDisposition {
	case dokan.FileSupersede:
		// FileSupersede   = CreateDisposition(0) If the file already exists, replace
		//it with the given file. If it does not, create the given file.

	case dokan.FileOpen:
		// FileOpen        = CreateDisposition(1) If the file already exists, open it
		//instead of creating a new file. If it does not, fail the request and do
		//not create a new file
		var f *BoltFile
		if err = fs.db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket(BucketNameMetadata)
			f, err = LoadBoltFile(b, fi.Path())
			if err != nil {
				if os.IsNotExist(err) {
					return dokan.ErrObjectPathNotFound //file doesnt exist
				}

				return err //other error
			}

			//file exists and deserialized
			return nil
		}); err != nil {
			return nil, false, err
		}

		return f, f.IsDir(), nil
	case dokan.FileCreate:
		// FileCreate      = CreateDisposition(2) If the file already exists, fail
		//the request and do not create or open the given file. If it does not,
		//create the given file.

	case dokan.FileOpenIf:
		// FileOpenIf      = CreateDisposition(3) If the file already exists, open
		//it. If it does not, create the given file.

	case dokan.FileOverwrite:
		// FileOverwrite   = CreateDisposition(4) If the file already exists, open
		//it and overwrite it. If it does not, fail the request.

	case dokan.FileOverwriteIf:
		// FileOverwriteIf = CreateDisposition(5) If the file already exists, open
		//it and overwrite it. If it does not, create the given file.
		var f *BoltFile
		if err = fs.db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket(BucketNameMetadata)
			f = NewBoltFile(false)
			//@TODO set populate file attributes

			//always overwrite file in db
			err = f.Save(b, fi.Path())
			if err != nil {
				return err
			}

			return nil
		}); err != nil {
			return nil, false, err
		}

		return f, f.IsDir(), nil
	}

	return nil, false, dokan.ErrNotSupported
}
