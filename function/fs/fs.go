package fs

import (
	"io"
	"os"
)

type ReaderFunc func(p []byte) (n int, err error)
type WriterFunc func(p []byte) (n int, err error)
type CloserFunc func() error

type FilesystemFunc func(name string) (os.File, error)
type WalkFunc func(path string, info os.FileInfo, err error) error
type WalkerFunc func(root string, walkFn WalkFunc) error
type RenderFunc func(wr io.Writer, name string, data interface{}) error
