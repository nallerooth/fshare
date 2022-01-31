package server

import "fmt"

type File struct {
	Filename  string
	Size      int64
	CreatedAt int64
}

func (f *File) String() string {
	return fmt.Sprintf("%s : %s", f.Filename, f.HumanSize())
}

func (f *File) HumanSize() string {
	units := []string{"", "K", "M", "G"}

	s := float64(f.Size)
	unit := 0
	for ; s > 1024 && unit < len(units); unit++ {
		s /= 1024
	}

	return fmt.Sprintf("%.2f%s", s, units[unit])
}
