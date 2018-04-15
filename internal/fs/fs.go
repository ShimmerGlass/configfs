package fs

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/process"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
)

func Run(mountPoint string, cb func(path string) ([]byte, error)) error {
	nfs := pathfs.NewPathNodeFs(&fs{
		FileSystem: pathfs.NewDefaultFileSystem(),
		cb:         cb,
	}, nil)
	server, _, err := nodefs.MountRoot(mountPoint, nfs.Root(), nil)
	if err != nil {
		return err
	}
	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChannel
		log.Println("unmounting")
		if err := server.Unmount(); err != nil {
			log.Println(err)
		}
		os.Exit(0)
	}()
	log.Println("mounting on", mountPoint)
	server.Serve()
	time.Sleep(time.Second)
	return nil
}

type fs struct {
	pathfs.FileSystem
	cb func(path string) ([]byte, error)
}

func (me *fs) GetAttr(name string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	switch name {
	case "":
		return &fuse.Attr{
			Mode: fuse.S_IFDIR | 0755,
		}, fuse.OK
	default:
		contents, err := me.contents(context)
		if err != nil {
			log.Println(err)
			return nil, fuse.ENODATA
		}

		return &fuse.Attr{
			Mode: fuse.S_IFREG | 0644, Size: uint64(len(contents)),
		}, fuse.OK
	}
}

func (me *fs) OpenDir(name string, context *fuse.Context) (c []fuse.DirEntry, code fuse.Status) {
	return []fuse.DirEntry{}, fuse.OK
}

func (fs *fs) Access(name string, mode uint32, context *fuse.Context) (code fuse.Status) {
	return fuse.OK
}

func (me *fs) Open(name string, flags uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	contents, err := me.contents(context)
	if err != nil {
		log.Println(err)
		return nil, fuse.ENODATA
	}

	return nodefs.NewDataFile(contents), fuse.OK
}

func (fs *fs) contents(context *fuse.Context) ([]byte, error) {
	path, err := fs.path(context)
	if err != nil {
		return nil, err
	}

	return fs.cb(path)
}

func (fs *fs) path(context *fuse.Context) (string, error) {
	proc, err := process.NewProcess(int32(context.Pid))
	if err != nil {
		return "", errors.Wrap(err, "caller process open")
	}

	wd, err := proc.Cwd()
	if err != nil {
		return "", errors.Wrap(err, "caller process cwd")
	}
	ex, err := proc.Exe()
	if err != nil {
		return "", errors.Wrap(err, "caller process ex")
	}

	log.Println("access from", ex, "in", wd)

	return wd, nil
}
