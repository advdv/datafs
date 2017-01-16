package datafs

import (
	"crypto/sha1"
	"errors"
)

var (
	//ErrExists is returned when a file or directory already exists when creating
	ErrExists = errors.New("File or directory already exists")

	//ErrNotExist is returned when there is no file while it is expected
	ErrNotExist = errors.New("No such file or directory")

	//ErrNotDirectory is returned when the file is not a directory while it is expected to
	ErrNotDirectory = errors.New("Not a directory")
)

var (
	//BucketNameFiles refers to the bucket that holds file (metadata)
	BucketNameFiles = []byte("files")

	//BucketNameChunks refers to the bucket that holds file contents
	BucketNameChunks = []byte("chunks")
)

//FileSystem maps file system semantics unto the bolt db buckets that
//is abstract enough that it can be used by OS specific user
//land file system proxies (FUSE, Dokany):
// - On windows it was designed for implementing Dokany's:
//     * CreateFile(ctx context.Context, fi *FileInfo, data *CreateData) (file File, isDirectory bool, err error)
//       to open, create and overwrite files or directories
//     * MoveFile(ctx context.Context, source *FileInfo, targetPath string, replaceExisting bool) error
// - On Linux (or OSX) FUSE it is modelled around the requests of
//     * Open(ctx context.Context, req *fuse.OpenRequest, resp *fuse.OpenResponse) (fs.Handle, error)
//     * Attr(ctx context.Context, a *fuse.Attr) error
//     * ReadDirAll(ctx context.Context) ([]fuse.Dirent, error)
//     * Write(ctx context.Context, req *fuse.WriteRequest, resp *fuse.WriteResponse) error
type FileSystem struct{}

//File are hold the metadata information for a path in the fileystem
//tree. It may be a directory (under the prefix of some other files)
//or reference a list of chunks that can be streamd as file content
// - On windows it was modelled for Dokany's
//      * WriteFile(ctx context.Context, fi *FileInfo, bs []byte, offset int64) (int, error)
//      * ReadFile(ctx context.Context, fi *FileInfo, bs []byte, offset int64) (int, error)
// - On Linux (or OSX) Fuse it is modelled for:
//     * Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error
//     * Write(ctx context.Context, req *fuse.WriteRequest, resp *fuse.WriteResponse) error
type File struct{}

//Read will get bytes from a file's chunked content and place then into buffer 'buf'
func (f *File) Read(offset int64, buf []byte) (n int, err error) {
	return 0, nil
}

//Write will put bytes a file's chunked content from buffer 'buf'
//@TODO how do we deal with buffering before chunking chunking?
func (f *File) Write(offset int64, buf []byte) (n int, err error) {
	return 0, nil
}

//K is the fixd-sized hash key of a piece of file content
type K [sha1.Size]byte

//Chunk holds a arbitrary-sized piece of file content
type Chunk []byte

//Open returns a file at joined path p, depending on provied
//flags it may create necessary files (or not) on demand.
func (fs *FileSystem) Open(p ...[]byte) (f *File, err error) {

	return nil, nil
}

//List reads files at jained path elements 'p', if it doesn't
//refer to a directory, an ErrNotDirectory is returned
func (fs *FileSystem) List(p ...[]byte) (ls []*File, err error) {
	return ls, nil
}
