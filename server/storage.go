package server

import (
	"crypto/sha256"
	"fmt"
	"os"
)

// hashFileName returns a hash of the salted filename
func (s *Server) hashFileName(filename string) string {
	saltedName := filename + s.config.Salt
	return fmt.Sprintf("%x", sha256.Sum256([]byte(saltedName)))
}

func (s *Server) findFile(files HashFileMap, hash string) (*File, error) {
	filename, ok := files[hash]
	if ok {
		return filename, nil
	}

	return nil, HashNotFoundError{Hash: hash}
}

// loadWorkdir loads the files in the specified workdir and calculates hashes
// for them, allowing fetching files without using the real file name
func (s *Server) LoadWorkdir() error {
	files, err := os.ReadDir(s.config.Workdir)
	if err != nil {
		return fmt.Errorf("Error: Unable to verify workdir -> %s", err)
	}

	mapped := HashFileMap{}

	// iterate over workdir contents, ignoring directories
	for _, f := range files {
		if !f.IsDir() {
			fileInfo, err := f.Info()
			if err != nil {
				return err
			}
			mapped[s.hashFileName(f.Name())] = &File{
				Filename: string(f.Name()),
				Size:     fileInfo.Size(),
			}
		}
	}

	s.files = mapped

	return nil
}

func (s *Server) AvailableFiles() HashFileMap {
	return s.files
}
