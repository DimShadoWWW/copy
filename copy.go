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
func Copy(c *slurp.C, keepath bool) slurp.Stage {
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

			s, err := file.Stat()
			if err != nil {
				c.Error(err)
				break
			}
			name := s.Name()
			file.Dir = ""

			if keepath {
				name = strings.Join(strings.Split(file.Path, "/")[1:], "/")
				file.Dir = strings.Split(file.Dir, "/")[0]
			}
			f := NewFile(name, file.Dir, buf.Bytes())

			fs = append(fs, f)

			file.Path = name
			file.Reader = buf
			file.FileInfo.SetSize(int64(buf.Len()))

			out <- file
		}
	}
}
