package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"os"

	"github.com/nallerooth/fshare/server"
)

var config = server.Config{
	Port:       32000,
	Passphrase: "",
	Salt:       "<change_me>",
	URL:        "https://nallerooth.com/share/",
	Workdir:    ".",
}

func init() {
	flag.UintVar(&config.Port, "port", config.Port, "Listening port")
	flag.StringVar(&config.Passphrase, "pass", config.Passphrase, "Passphrase for uploading files [optional]")
	flag.StringVar(&config.Salt, "salt", config.Salt, "Salt")
	flag.StringVar(&config.Workdir, "wd", config.Workdir, "Workdir")

	flag.Parse()
}

// hashFileName returns a hash of the salted filename
func hashFileName(fn string) string {
	saltedName := fn + config.Salt
	return fmt.Sprintf("%x", sha256.Sum256([]byte(saltedName)))
}

func findFile(files server.HashFileMap, hash string) (*server.File, error) {
	filename, ok := files[hash]
	if ok {
		return filename, nil
	}

	return nil, server.HashNotFoundError{hash}
}

// loadWorkdir loads the files in the specified workdir and calculates hashes
// for them, allowing fetching files without using the real file name
func loadWorkdir() (server.HashFileMap, error) {
	files, err := os.ReadDir(config.Workdir)
	if err != nil {
		return nil, fmt.Errorf("Error: Unable to verify workdir -> %s", err)
	}

	lf := server.HashFileMap{}

	// iterate over workdir contents, ignoring directories
	for _, f := range files {
		if !f.IsDir() {
			fileInfo, err := f.Info()
			if err != nil {
				return nil, err
			}
			lf[hashFileName(f.Name())] = &server.File{
				Filename: string(f.Name()),
				Size:     fileInfo.Size(),
			}
		}
	}

	return lf, nil
}

func main() {
	fmt.Printf("%+v\n", config)

	storage, err := loadWorkdir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading workdir: %s\n", err)
	}

	for hash, name := range storage {
		fmt.Println(hash, " -> ", name)
	}
}
