package configfs

import (
	"github.com/prometheus/common/log"
	"golang.org/x/net/context"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
)

type contentFn func() ([]byte, error)
type doneFn func()

type FS struct {
	fn contentFn
	fuse *fuse.Conn
}

func MountFS(mountPoint string, contentFn contentFn) (doneFn, <-chan error) {
	out := make(chan error, 1)

	c, err := fuse.Mount(
		mountPoint,
		fuse.FSName("configfs"),
		fuse.Subtype("configfs"),
		fuse.LocalVolume(),
		fuse.VolumeName("Config"),
	)
	if err != nil {
		out <- err
	}

	go func() {
		log.Info("Serving...")
		err = fs.Serve(c, FS{
			fn: contentFn,
		})
		if err != nil {
			out <- err
		}
	}()

	go func() {
		<-c.Ready
		log.Info("Ready")
		if err := c.MountError; err != nil {
			out <- err
		}
	}()

	return func() {
		log.Info("Unmounting")
		err := fuse.Unmount(mountPoint)
		if err != nil {
			out <- err
		}

		err = c.Close()
		if err != nil {
			out <- err
		}
	}, out
}

func (f FS) Root() (fs.Node, error) {
	return &File{
		fn: f.fn,
	}, nil
}

type File struct {
	fn contentFn
}

func (f File) Attr(ctx context.Context, a *fuse.Attr) error {
	contents, err := f.fn()
	if err != nil {
		log.Error(err)
		return err
	}

	a.Inode = 2
	a.Mode = 0444
	a.Size = uint64(len(contents))
	return nil
}

func (f File) ReadAll(ctx context.Context) ([]byte, error) {
	contents, err := f.fn()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return []byte(contents), nil
}
