package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/keybase/kbfs/dokan"
	"github.com/keybase/kbfs/dokan/winacl"
	"golang.org/x/net/context"
)

func debug(args ...interface{}) {
	log.Println(args...)
}

func debugf(s string, args ...interface{}) {
	log.Printf(s, args...)
}

type emptyFS struct{}

func (t emptyFile) GetFileSecurity(ctx context.Context, fi *dokan.FileInfo, si winacl.SecurityInformation, sd *winacl.SecurityDescriptor) error {
	debug("emptyFS.GetFileSecurity")
	return nil
}
func (t emptyFile) SetFileSecurity(ctx context.Context, fi *dokan.FileInfo, si winacl.SecurityInformation, sd *winacl.SecurityDescriptor) error {
	debug("emptyFS.SetFileSecurity")
	return nil
}
func (t emptyFile) Cleanup(ctx context.Context, fi *dokan.FileInfo) {
	debug("emptyFS.Cleanup")
}

func (t emptyFile) CloseFile(ctx context.Context, fi *dokan.FileInfo) {
	debug("emptyFS.CloseFile")
}

func (t emptyFS) WithContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return ctx, nil
}

func (t emptyFS) GetVolumeInformation(ctx context.Context) (dokan.VolumeInformation, error) {
	debug("emptyFS.GetVolumeInformation")
	return dokan.VolumeInformation{}, nil
}

func (t emptyFS) GetDiskFreeSpace(ctx context.Context) (dokan.FreeSpace, error) {
	debug("emptyFS.GetDiskFreeSpace")
	return dokan.FreeSpace{}, nil
}

func (t emptyFS) ErrorPrint(err error) {
	debug(err)
}

func (t emptyFS) CreateFile(ctx context.Context, fi *dokan.FileInfo, cd *dokan.CreateData) (dokan.File, bool, error) {
	debug("emptyFS.CreateFile")
	return emptyFile{}, true, nil
}
func (t emptyFile) CanDeleteFile(ctx context.Context, fi *dokan.FileInfo) error {
	return dokan.ErrAccessDenied
}
func (t emptyFile) CanDeleteDirectory(ctx context.Context, fi *dokan.FileInfo) error {
	return dokan.ErrAccessDenied
}
func (t emptyFile) SetEndOfFile(ctx context.Context, fi *dokan.FileInfo, length int64) error {
	debug("emptyFile.SetEndOfFile")
	return nil
}
func (t emptyFile) SetAllocationSize(ctx context.Context, fi *dokan.FileInfo, length int64) error {
	debug("emptyFile.SetAllocationSize")
	return nil
}
func (t emptyFS) MoveFile(ctx context.Context, source *dokan.FileInfo, targetPath string, replaceExisting bool) error {
	debug("emptyFS.MoveFile")
	return nil
}
func (t emptyFile) ReadFile(ctx context.Context, fi *dokan.FileInfo, bs []byte, offset int64) (int, error) {
	return len(bs), nil
}
func (t emptyFile) WriteFile(ctx context.Context, fi *dokan.FileInfo, bs []byte, offset int64) (int, error) {
	return len(bs), nil
}
func (t emptyFile) FlushFileBuffers(ctx context.Context, fi *dokan.FileInfo) error {
	debug("emptyFS.FlushFileBuffers")
	return nil
}

type emptyFile struct{}

func (t emptyFile) GetFileInformation(ctx context.Context, fi *dokan.FileInfo) (*dokan.Stat, error) {
	debug("emptyFile.GetFileInformation")
	var st dokan.Stat
	st.FileAttributes = dokan.FileAttributeNormal
	return &st, nil
}
func (t emptyFile) FindFiles(context.Context, *dokan.FileInfo, string, func(*dokan.NamedStat) error) error {
	debug("emptyFile.FindFiles")
	return nil
}
func (t emptyFile) SetFileTime(context.Context, *dokan.FileInfo, time.Time, time.Time, time.Time) error {
	debug("emptyFile.SetFileTime")
	return nil
}
func (t emptyFile) SetFileAttributes(ctx context.Context, fi *dokan.FileInfo, fileAttributes dokan.FileAttribute) error {
	debug("emptyFile.SetFileAttributes")
	return nil
}

func (t emptyFile) LockFile(ctx context.Context, fi *dokan.FileInfo, offset int64, length int64) error {
	debug("emptyFile.LockFile")
	return nil
}
func (t emptyFile) UnlockFile(ctx context.Context, fi *dokan.FileInfo, offset int64, length int64) error {
	debug("emptyFile.UnlockFile")
	return nil
}

func main() {
	log.Printf("started")
	defer log.Printf("exited")
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	fs := &emptyFS{}
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
