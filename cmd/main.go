package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/mackee/quasar"
)

func main() {
	var filename string
	flag.StringVar(&filename, "config", "quasar.yml", "configuration file")
	flag.Parse()

	c, err := quasar.ParseConfig(filename)
	if err != nil {
		log.Fatal(err)
	}

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println(`
Sub Commands:
serve  : serving daemon and RPC server.
GetEnv : retrieve enviroment string from RPC server.
Close  : close enviroment signal to RPC server.
		`)
		return
	}

	switch args[0] {
	case "serve":
		quasar.Run(c)
	case "GetEnv", "Close":
		if len(args) < 3 {
			log.Fatalf(
				"not enough args.\nExample: quasar %s <DaemoName> <Envname>\n",
				args[0],
			)
		}
		client, err := quasar.NewClient(c)
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()
		switch args[0] {
		case "GetEnv":
			resp, err := client.GetEnv(args[1], args[2])
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(resp)
		case "Close":
			err := client.EnvClose(args[1], args[2])
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
