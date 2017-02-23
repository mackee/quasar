package quasar

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
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
	inss := make([]instance, 0, len(config.Daemons))
	for _, d := range config.Daemons {
		ins, err := d.ToInstance()
		if err != nil {
			log.Printf("fail instance %s: %s", d.Name, err)
			continue
		}
		inss = append(inss, ins)
		wg.Add(1)
		go func(ins instance) {
			defer wg.Done()
			err = ins.Run(ctx)
			if err != nil {
				log.Printf("cannot start %s: %s", d.Name, err)
				return
			}
			ins.Wait()
		}(ins)
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		for _, ins := range inss {
			go ins.Stop()
		}
	}()
	wg.Wait()
}
