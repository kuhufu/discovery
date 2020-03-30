package main

import (
	"github.com/kuhufu/discovery"
	"log"
	"testing"
	"time"
)

func TestRegister(t *testing.T) {
	register("loclahost:8070")
	register("loclahost:8080")
	register("loclahost:8090")

	time.Sleep(time.Hour)
}

func register(addr string) {
	r := discovery.NewRegister(
		discovery.Name("srv.foo"),
		discovery.RegisterAddrs("localhost:2379"),
		discovery.Version("v2"),
	)

	err := r.Register(addr)
	if err != nil {
		log.Print(err)
		return
	}
}

func TestResolver(t *testing.T) {
	r := discovery.NewResolver(
		discovery.Name("srv.foo"),
		discovery.RegisterAddrs("localhost:2379"),
		discovery.Version("v1"),
	)

	r.Watch()
}
