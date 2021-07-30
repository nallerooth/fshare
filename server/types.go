package server

import "fmt"

// ERRORS BELOW

type FileNotFoundError struct {
	Filename string
}

func (e FileNotFoundError) Error() string {
	return fmt.Sprintf("File not found: %s", e.Filename)
}

type HashNotFoundError struct {
	Hash string
}

func (e HashNotFoundError) Error() string {
	return fmt.Sprintf("Hash not found: %s", e.Hash)
}
