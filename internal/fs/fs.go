package fs

import (
	"encoding/base64"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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
		log.Println("Unmounting")
		log.Fatal("end", server.Unmount())
	}()
	server.Serve()
	time.Sleep(time.Second)
	return nil
}

type fs struct {
	pathfs.FileSystem
	cb func(path string) ([]byte, error)
}

func (me *fs) GetAttr(name string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	return &fuse.Attr{
		Mode: fuse.S_IFREG | 0644, Size: uint64(len(name)),
	}, fuse.OK
}

func (me *fs) OpenDir(name string, context *fuse.Context) (c []fuse.DirEntry, code fuse.Status) {
	return nil, fuse.ENOENT
}

func (me *fs) Open(name string, flags uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	log.Println("fuse file access", name)
	path, err := base64.URLEncoding.DecodeString(name)
	if err != nil {
		log.Println(err)
		return nil, fuse.ENOENT
	}

	contents, err := me.cb(string(path))
	if err != nil {
		log.Println(err)
		return nil, fuse.ENOENT
	}

	return nodefs.NewDataFile(contents), fuse.OK
}
