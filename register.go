package discovery

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"log"
	"time"
)

type Register struct {
	options Options
	cli     *clientv3.Client
	service Service
	exitC   chan struct{}
}

func NewRegister(opts ...Option) *Register {
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

	return &Register{
		options: options,
		cli:     cli,
		service: service,
		exitC:   make(chan struct{}),
	}
}

func (r *Register) Register(addr string) error {
	cli := r.cli
	leaseRsp, err := cli.Grant(context.Background(), 18)
	if err != nil {
		return err
	}

	node := &Node{
		Id:   fmt.Sprintf("%v", time.Now().UnixNano()),
		Addr: addr,
		Meta: nil,
	}

	r.service.Name = r.options.ServiceName
	r.service.Version = r.options.ServiceVersion
	r.service.nodes = append(r.service.nodes, )

	leaseId := leaseRsp.ID

	_, err = cli.Put(context.Background(), nodePath(r.service.Name, r.service.Version, addr), node.String(), clientv3.WithLease(leaseId))
	if err != nil {
		return err
	}

	leaseRspChan, err := cli.KeepAlive(context.Background(), leaseId)
	if err != nil {
		log.Printf("lease %v keepalive err: %v", leaseId, err)
		return err
	}

	go func() {
		for {
			select {
			case <-r.exitC:
				log.Printf("register exit")
				return
			case rsp := <-leaseRspChan:
				if rsp != nil {
					log.Printf("续租成功 %v", addr)
				} else {
					log.Printf("续租已关闭%v", addr)
					if err := r.Register(addr); err != nil {
						log.Printf("开启新租约失败%v", addr)
					} else {
						log.Printf("开启新租约成功%v", addr)
					}

					return
				}

			}
		}
	}()

	return nil
}

func (r *Register) Close() error {
	close(r.exitC)
	return nil
}
