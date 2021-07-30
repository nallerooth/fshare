package main

import (
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
	Workdir:    "./data/",
}

func init() {
	flag.UintVar(&config.Port, "port", config.Port, "Listening port")
	flag.StringVar(&config.Passphrase, "pass", config.Passphrase, "Passphrase for uploading files [optional]")
	flag.StringVar(&config.Salt, "salt", config.Salt, "Salt")
	flag.StringVar(&config.Workdir, "wd", config.Workdir, "Workdir")

	flag.Parse()
}
func main() {
	fmt.Printf("%+v\n\n", config)

	s, err := server.New(config)
	if err != nil {
		panic(err)
	}

	err = s.LoadWorkdir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading workdir: %s\n", err)
	}

	for hash, name := range s.AvailableFiles() {
		fmt.Println(hash, " -> ", name)
	}

	s.Listen()
}
