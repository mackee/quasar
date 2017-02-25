package quasar

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
)

func Run(c config) {
	ctx := context.Background()
	var wg sync.WaitGroup
	inss := map[string]instance{}
	for _, d := range c.Daemons {
		ins, err := d.ToInstance()
		if err != nil {
			log.Printf("fail instance %s: %s", d.Name, err)
			continue
		}
		inss[d.Name] = ins
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
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	go func() {
		<-sigCh
		for _, ins := range inss {
			go ins.Stop()
		}
	}()

	go Serve(c, inss)

	wg.Wait()
}
