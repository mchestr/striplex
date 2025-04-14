package main

import (
	"flag"
	"fmt"
	"os"

	"striplex/config"
	"striplex/db"
	"striplex/server"

	"github.com/stripe/stripe-go/v82"
)

func main() {
	environment := flag.String("e", "development", "")
	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode}")
		os.Exit(1)
	}
	flag.Parse()
	config.Init(*environment)

	fmt.Println(config.Config.GetStringSlice("plex.shared_libraries"))

	stripe.Key = config.Config.GetString("stripe.secret_key")
	db.Connect()
	server.Init()
}
