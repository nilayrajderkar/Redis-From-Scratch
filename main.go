package main

import (
	"github.com/nilayrajderkar/redis-implementation/client"
	"log"
)

func main() {
	if err := client.StartClient(); err != nil {
		log.Fatal(err)
	}
}
