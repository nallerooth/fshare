package server

import "fmt"

type File struct {
	Filename string
	Size     int64
	// TODO: Add created at
}

func (f *File) String() string {
	return fmt.Sprintf("%s : %s", f.Filename, f.HumanSize())
}

func (f *File) HumanSize() string {
	units := []string{"B", "KB", "MB", "GB"}

	s := float64(f.Size)
	unit := 0
	for ; s > 1024 && unit < len(units); unit++ {
		s /= 1024
	}

	return fmt.Sprintf("%.2f %s", s, units[unit])
}
