package server

import "fmt"

type File struct {
	Filename string
	Size     int64
}

func (f *File) String() string {
	return fmt.Sprintf("%s : %d", f.Filename, f.Size)
}
