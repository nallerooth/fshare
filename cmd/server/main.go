package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/nallerooth/fshare/internal/config"
	"github.com/nallerooth/fshare/internal/server"
)

// Default config, will be overridden by flags
var conf = config.Config{
	Server: config.ServerConfig{
		Port:      32000,
		PublicURL: "https://nallerooth.com/fshare/",
		Workdir:   "",
		Salt:      "<change_me>",
	},
}

func init() {
	flag.UintVar(&conf.Server.Port, "port", conf.Server.Port, "Listening port")
	flag.StringVar(&conf.Passphrase, "pass", conf.Passphrase, "Passphrase for uploading files [optional]")
	flag.StringVar(&conf.Server.Salt, "salt", conf.Server.Salt, "Salt")
	flag.StringVar(&conf.Server.Workdir, "wd", conf.Server.Workdir, "Workdir")
	flag.StringVar(&conf.Server.PublicURL, "url", conf.Server.PublicURL, "Public URL to reach shared files")

	flag.Parse()
}

func main() {
	fmt.Printf("%+v\n\n", conf)

	// TODO: verify that a valid config has been loaded (either flags or conf file)

	s, err := server.NewServer(conf)
	if err != nil {
		panic(err)
	}

	err = s.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %s\n", err)
	}
}
