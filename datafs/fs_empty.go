package datafs

import (
	"time"

	"github.com/keybase/kbfs/dokan"
	"github.com/keybase/kbfs/dokan/winacl"
	"golang.org/x/net/context"
)

//EmptyFS is a noop dokan filesystem
type EmptyFS struct{}

// WithContext returns a context for a new request. If the CancelFunc
// is not null, it is called after the request is done. The most minimal
// implementation is
// `func (*T)WithContext(c context.Context) { return c, nil }`.
func (t EmptyFS) WithContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return ctx, nil
}

// GetVolumeInformation returns information about the volume.
func (t EmptyFS) GetVolumeInformation(ctx context.Context) (dokan.VolumeInformation, error) {
	debug("EmptyFS.GetVolumeInformation")
	return dokan.VolumeInformation{}, nil
}

// GetDiskFreeSpace returns information about disk free space.
// Called quite often by Explorer.
func (t EmptyFS) GetDiskFreeSpace(ctx context.Context) (dokan.FreeSpace, error) {
	debug("EmptyFS.GetDiskFreeSpace")
	return dokan.FreeSpace{}, nil
}

// ErrorPrint is called when dokan needs notify the program of an error message.
// A sensible approach is to print the error.
func (t EmptyFS) ErrorPrint(err error) {
	debug(err)
}

// CreateFile is called to open and create files.
func (t EmptyFS) CreateFile(ctx context.Context, fi *dokan.FileInfo, cd *dokan.CreateData) (dokan.File, bool, error) {
	debug("EmptyFS.CreateFile(ctx, '"+fi.Path(), "', cd")
	return EmptyFile{}, true, nil
}

// MoveFile corresponds to rename.
func (t EmptyFS) MoveFile(ctx context.Context, source *dokan.FileInfo, targetPath string, replaceExisting bool) error {
	debug("EmptyFS.MoveFile")
	return nil
}

// EmptyFile is the noop interface for files and directories.
type EmptyFile struct{}

//GetFileSecurity gets specified information about the security of a file or directory.
func (t EmptyFile) GetFileSecurity(ctx context.Context, fi *dokan.FileInfo, si winacl.SecurityInformation, sd *winacl.SecurityDescriptor) error {
	debug("EmptyFS.GetFileSecurity")
	return nil
}

//SetFileSecurity sets the security of a file or directory object.
func (t EmptyFile) SetFileSecurity(ctx context.Context, fi *dokan.FileInfo, si winacl.SecurityInformation, sd *winacl.SecurityDescriptor) error {
	debug("EmptyFS.SetFileSecurity")
	return nil
}

// Cleanup is called after the last handle from userspace is closed.
// Cleanup must perform actual deletions marked from CanDelete*
// by checking FileInfo.IsDeleteOnClose if the filesystem supports
// deletions.
func (t EmptyFile) Cleanup(ctx context.Context, fi *dokan.FileInfo) {
	debug("EmptyFS.Cleanup")
}

// CloseFile is called when closing a handle to the file.
func (t EmptyFile) CloseFile(ctx context.Context, fi *dokan.FileInfo) {
	debug("EmptyFS.CloseFile")
}

// CanDeleteFile and CanDeleteDirectory should check whether the file/directory
// can be deleted. The actual deletion should be done by checking
// FileInfo.IsDeleteOnClose in Cleanup.
func (t EmptyFile) CanDeleteFile(ctx context.Context, fi *dokan.FileInfo) error {
	return dokan.ErrAccessDenied
}

// CanDeleteDirectory should check whether the file/directory
// can be deleted. The actual deletion should be done by checking
// FileInfo.IsDeleteOnClose in Cleanup.
func (t EmptyFile) CanDeleteDirectory(ctx context.Context, fi *dokan.FileInfo) error {
	return dokan.ErrAccessDenied
}

// SetEndOfFile truncates the file. May be used to extend a file with zeros.
func (t EmptyFile) SetEndOfFile(ctx context.Context, fi *dokan.FileInfo, length int64) error {
	debug("EmptyFile.SetEndOfFile")
	return nil
}

// SetAllocationSize see FILE_ALLOCATION_INFORMATION on MSDN.
// For simple semantics if length > filesize then ignore else truncate(length).
func (t EmptyFile) SetAllocationSize(ctx context.Context, fi *dokan.FileInfo, length int64) error {
	debug("EmptyFile.SetAllocationSize")
	return nil
}

// ReadFile implements read for dokan.
func (t EmptyFile) ReadFile(ctx context.Context, fi *dokan.FileInfo, bs []byte, offset int64) (int, error) {
	return len(bs), nil
}

// WriteFile implements write for dokan.
func (t EmptyFile) WriteFile(ctx context.Context, fi *dokan.FileInfo, bs []byte, offset int64) (int, error) {
	return len(bs), nil
}

// FlushFileBuffers corresponds to fsync.
func (t EmptyFile) FlushFileBuffers(ctx context.Context, fi *dokan.FileInfo) error {
	debug("EmptyFS.FlushFileBuffers")
	return nil
}

// GetFileInformation - corresponds to stat.
func (t EmptyFile) GetFileInformation(ctx context.Context, fi *dokan.FileInfo) (*dokan.Stat, error) {
	debug("EmptyFile.GetFileInformation(ctx, '" + fi.Path() + "')")
	var st dokan.Stat
	st.FileAttributes = dokan.FileAttributeNormal
	return &st, nil
}

// FindFiles is the readdir. The function is a callback that should be called
// with each file. The same NamedStat may be reused for subsequent calls.
//
// Pattern will be an empty string unless UseFindFilesWithPattern is enabled - then
// it may be a pattern like `*.png` to match. All implementations must be prepared
// to handle empty strings as patterns.
func (t EmptyFile) FindFiles(ctx context.Context, fi *dokan.FileInfo, pattern string, fn func(*dokan.NamedStat) error) error {
	debug("EmptyFile.FindFiles(ctx, '" + fi.Path() + "', '" + pattern + "', fn")
	return nil
}

// SetFileTime sets the file time. Test times with .IsZero
// whether they should be set.
func (t EmptyFile) SetFileTime(context.Context, *dokan.FileInfo, time.Time, time.Time, time.Time) error {
	debug("EmptyFile.SetFileTime")
	return nil
}

// SetFileAttributes is for setting file attributes.
func (t EmptyFile) SetFileAttributes(ctx context.Context, fi *dokan.FileInfo, fileAttributes dokan.FileAttribute) error {
	debug("EmptyFile.SetFileAttributes")
	return nil
}

// LockFile at a specific offset and data length. This is only used if \ref DOKAN_OPTION_FILELOCK_USER_MODE is enabled.
func (t EmptyFile) LockFile(ctx context.Context, fi *dokan.FileInfo, offset int64, length int64) error {
	debug("EmptyFile.LockFile")
	return nil
}

//UnlockFile at a specific offset and data lengthh. This is only used if \ref DOKAN_OPTION_FILELOCK_USER_MODE is enabled.
func (t EmptyFile) UnlockFile(ctx context.Context, fi *dokan.FileInfo, offset int64, length int64) error {
	debug("EmptyFile.UnlockFile")
	return nil
}
