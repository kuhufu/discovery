package discovery

import (
	"context"
	"go.etcd.io/etcd/clientv3"
	"log"
)

type Resolver struct {
	options Options
	cli     *clientv3.Client
	service Service
	exitC   chan struct{}
}

func NewResolver(opts ...Option) *Resolver {
	options := newOptions(opts...)

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   options.Addrs,
		DialTimeout: options.DialTimeout,
	})
	if err != nil {
		panic(err)
	}

	service := Service{
		Name:    options.ServiceName,
		Version: options.ServiceVersion,
	}

	return &Resolver{
		options: options,
		cli:     cli,
		service: service,
		exitC:   make(chan struct{}),
	}
}

func (r *Resolver) fullSync() error {
	resp, err := r.cli.Get(context.Background(), nodePrefix(r.service.Name, r.service.Version), clientv3.WithPrefix())
	if err != nil {
		return err
	}

	var nodes []*Node
	for _, kv := range resp.Kvs {
		node := &Node{}
		err := node.Decode(kv.Value)
		if err != nil {
			return err
		}
		nodes = append(nodes, node)
	}

	r.service.nodes = nodes

	return nil
}

func (r *Resolver) Watch() error {
	err := r.fullSync()
	if err != nil {
		return err
	}

	cli := r.cli

	watchChan := cli.Watch(context.Background(), nodePrefix(r.service.Name, r.service.Version), clientv3.WithPrefix())

	for {
		select {
		case <-r.exitC:
			log.Printf("resolver exit")
			return nil
		case rsp := <-watchChan:
			for _, event := range rsp.Events {
				log.Printf("type:%v, key: %s, val: %s", event.Type, event.Kv.Key, event.Kv.Value)
				switch event.Type {
				case clientv3.EventTypeDelete:
				case clientv3.EventTypePut:

				}
			}
		}
	}
}
