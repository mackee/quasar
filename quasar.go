package quasar

import (
	"context"
	"flag"
	"log"
	"sync"
)

func Run() {
	var filename string
	flag.StringVar(&filename, "config", "quasar.yml", "configuration file")

	flag.Parse()

	config, err := ParseConfig(filename)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	var wg sync.WaitGroup
	for _, d := range config.Daemons {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ins, err := d.ToInstance()
			if err != nil {
				log.Printf("fail instance %s: %s", d.Name, err)
				return
			}
			err = ins.Run(ctx)
			if err != nil {
				log.Printf("cannot start %s: %s", d.Name, err)
				return
			}
			ins.Wait()
		}()
	}
	wg.Wait()
}
