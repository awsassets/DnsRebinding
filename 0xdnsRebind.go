package main

import (
	"0xdns-rebind/conf"
	"0xdns-rebind/core"
	"fmt"
	"log"
	"os"
	"os/signal"
)

func main() {
	addr := ":53"
	if err := conf.SetFromFile("config.yml"); err != nil {
		log.Fatalln(err)
	}
	dns, err := core.NewDNSDog(":53")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("0x-dns-rebind strat...")
	fmt.Printf("dns@'%s'\n", addr)

	go func() {
		if err := dns.ListenAndServe(); err != nil {
			log.Fatalln(err)
		}
	}()
	defer func() {
		if err := dns.Shutdown(); err != nil {
			log.Fatalln(err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
