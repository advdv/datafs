# datafs
A de-duplicated filesystem that syncs to S3 and works on OSX, Windows and Linux

## Anatomy of a Go -> C -> Dokany

kbfs/dokan.File interface

```Go
// File is the interface for files and directories.
type File interface {
	// ReadFile implements read for dokan.
	ReadFile(ctx context.Context, fi *FileInfo, bs []byte, offset int64) (int, error)
	// WriteFile implements write for dokan.
	WriteFile(ctx context.Context, fi *FileInfo, bs []byte, offset int64) (int, error)
	// FlushFileBuffers corresponds to fsync.
	FlushFileBuffers(ctx context.Context, fi *FileInfo) error

	// GetFileInformation - corresponds to stat.
	GetFileInformation(ctx context.Context, fi *FileInfo) (*Stat, error)

	// FindFiles is the readdir. The function is a callback that should be called
	// with each file. The same NamedStat may be reused for subsequent calls.
	//
	// Pattern will be an empty string unless UseFindFilesWithPattern is enabled - then
	// it may be a pattern like `*.png` to match. All implementations must be prepared
	// to handle empty strings as patterns.
	FindFiles(ctx context.Context, fi *FileInfo, pattern string, fillStatCallback func(*NamedStat) error) error

	// SetFileTime sets the file time. Test times with .IsZero
	// whether they should be set.
	SetFileTime(ctx context.Context, fi *FileInfo, creation time.Time, lastAccess time.Time, lastWrite time.Time) error
	// SetFileAttributes is for setting file attributes.
	SetFileAttributes(ctx context.Context, fi *FileInfo, fileAttributes FileAttribute) error

	// SetEndOfFile truncates the file. May be used to extend a file with zeros.
	SetEndOfFile(ctx context.Context, fi *FileInfo, length int64) error
	// SetAllocationSize see FILE_ALLOCATION_INFORMATION on MSDN.
	// For simple semantics if length > filesize then ignore else truncate(length).
	SetAllocationSize(ctx context.Context, fi *FileInfo, length int64) error

	LockFile(ctx context.Context, fi *FileInfo, offset int64, length int64) error
	UnlockFile(ctx context.Context, fi *FileInfo, offset int64, length int64) error

	GetFileSecurity(ctx context.Context, fi *FileInfo, si winacl.SecurityInformation, sd *winacl.SecurityDescriptor) error
	SetFileSecurity(ctx context.Context, fi *FileInfo, si winacl.SecurityInformation, sd *winacl.SecurityDescriptor) error

	// CanDeleteFile and CanDeleteDirectory should check whether the file/directory
	// can be deleted. The actual deletion should be done by checking
	// FileInfo.IsDeleteOnClose in Cleanup.
	CanDeleteFile(ctx context.Context, fi *FileInfo) error
	CanDeleteDirectory(ctx context.Context, fi *FileInfo) error
	// Cleanup is called after the last handle from userspace is closed.
	// Cleanup must perform actual deletions marked from CanDelete*
	// by checking FileInfo.IsDeleteOnClose if the filesystem supports
	// deletions.
	Cleanup(ctx context.Context, fi *FileInfo)
	// CloseFile is called when closing a handle to the file.
	CloseFile(ctx context.Context, fi *FileInfo)
}
```


kbfs/dokan.FileSystem interface

```Go
// FileSystem is the inteface for filesystems in Dokan.
type FileSystem interface {
	// WithContext returns a context for a new request. If the CancelFunc
	// is not null, it is called after the request is done. The most minimal
	// implementation is
	// `func (*T)WithContext(c context.Context) { return c, nil }`.
	WithContext(context.Context) (context.Context, context.CancelFunc)

	// CreateFile is called to open and create files.
	CreateFile(ctx context.Context, fi *FileInfo, data *CreateData) (file File, isDirectory bool, err error)

	// GetDiskFreeSpace returns information about disk free space.
	// Called quite often by Explorer.
	GetDiskFreeSpace(ctx context.Context) (FreeSpace, error)

	// GetVolumeInformation returns information about the volume.
	GetVolumeInformation(ctx context.Context) (VolumeInformation, error)

	// MoveFile corresponds to rename.
	MoveFile(ctx context.Context, source *FileInfo, targetPath string, replaceExisting bool) error

	// ErrorPrint is called when dokan needs notify the program of an error message.
	// A sensible approach is to print the error.
	ErrorPrint(error)
}
```

kbfs/dokan.Conf struct

```Go
type Config struct {
	// Path is the path to mount, e.g. `L:`. Must be set.
	Path string
	// FileSystem is the filesystem implementation. Must be set.
	FileSystem FileSystem
	// MountFlags for this filesystem instance. Is optional.
	MountFlags MountFlag
	// DllPath is the optional full path to dokan1.dll.
	// Empty causes dokan1.dll to be loaded from the system directory.
	// Only the first load of a dll determines the path -
	// further instances in the same process will use
	// the same instance regardless of path.
	DllPath string
}
```


kbfs/libdokan.Mounter interface

```Go
type Mounter interface {
	Dir() string
	Mount(*dokan.Config, logger.Logger) error
	Unmount() error
}
```

kbfs/libdokan.StartOptions struct
```Go
type StartOptions struct {
	KbfsParams  libkbfs.InitParams
	RuntimeDir  string
	Label       string
	DokanConfig dokan.Config
}
```

kbfs/libkbfs.Context interface
```Go
type Context interface {
	GetRunMode() libkb.RunMode
	GetLogDir() string
	GetDataDir() string
	ConfigureSocketInfo() (err error)
	GetSocket(clearError bool) (net.Conn, rpc.Transporter, bool, error)
	NewRPCLogFactory() *libkb.RPCLogFactory
}
```


then to start:
```
err = libdokan.Start(mounter, options, ctx)
```
