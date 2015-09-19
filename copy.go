package copy

// github.com/yosssi/gcss binding for slurp.
// No Configuration required.

import (
	"bytes"
	"strings"
	"sync"

	"github.com/omeid/slurp"
)

// File represents a file.
type File struct {
	path string
	dir  string
	data []byte
}

// NewFile creates and returns a file.
func NewFile(path, dir string, data []byte) *File {
	return &File{
		path: path,
		dir:  dir,
		data: data,
	}
}

// Copy the file content
func Copy(c *slurp.C) slurp.Stage {
	return func(in <-chan slurp.File, out chan<- slurp.File) {

		fs := []*File{}

		var wg sync.WaitGroup
		defer wg.Wait() //Wait before all templates are executed.

		for file := range in {
			buf := new(bytes.Buffer)
			_, err := buf.ReadFrom(file.Reader)
			file.Close()
			if err != nil {
				c.Error(err)
				continue
			}

			name := strings.Join(strings.Split(file.Path, "/")[1:], "/")
			dir := strings.Split(file.Dir, "/")[0]

			f := NewFile(name, dir, buf.Bytes())

			fs = append(fs, f)

			file.Dir = dir
			file.Reader = buf
			file.FileInfo.SetSize(int64(buf.Len()))

			out <- file

		}
	}
}
