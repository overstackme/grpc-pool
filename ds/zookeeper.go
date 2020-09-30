// +build zookeeper

package ds

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/go-zookeeper/zk"
	"google.golang.org/grpc/resolver"
)

type ZooKeeper struct {
	config  *Config
	client  *zk.Conn
	targets sync.Map
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewDiscovery(config *Config) Discovery {
	ctx, cancel := context.WithCancel(context.Background())
	return &ZooKeeper{
		config: config,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (z *ZooKeeper) Watcher() error {
	if err := z.connect(); err != nil {
		return err
	}
	go z.watcher()
	return nil
}

func (z *ZooKeeper) Registry(srvs []Service) error {
	if err := z.connect(); err != nil {
		return err
	}
	for _, srv := range srvs {
		if err := z.registry(srv); err != nil {
			return err
		}
	}
	return nil
}

func (z *ZooKeeper) Stop() {
	if z.client != nil {
		z.client.Close()
	}
}

func (z *ZooKeeper) Addresses(target resolver.Target) ([]resolver.Address, error) {
	value, ok := z.targets.Load(target.Endpoint)
	if !ok {
		z.targets.Store(target.Endpoint, nil)
	}
	targets, ok := value.([]string)
	if !ok {
		return nil, nil
	}

	var addrs []resolver.Address
	for _, v := range targets {
		addrs = append(addrs, resolver.Address{
			Addr: v,
		})
	}
	return addrs, nil
}

func (z *ZooKeeper) connect() error {
	client, _, err := zk.Connect(z.config.Endpoints, z.config.Timeout, zk.WithLogInfo(false))
	if err != nil {
		return err
	}
	z.client = client
	return nil
}

func (z *ZooKeeper) registry(srv Service) error {
	groupName := z.generalGroupName()
	fmt.Println(groupName)
	if err := z.createPathNode(groupName, 0); err != nil {
		return err
	}

	srvPath := fmt.Sprintf("%s/%s", groupName, srv.Name)
	fmt.Println(srvPath)
	if err := z.createPathNode(srvPath, 0); err != nil {
		return err
	}

	nodePath := fmt.Sprintf("%s/%s", srvPath, srv.Addr)
	fmt.Println(nodePath)
	if err := z.createPathNode(nodePath, zk.FlagEphemeral); err != nil {
		return err
	}
	return nil
}

func (z *ZooKeeper) createPathNode(path string, flag int) error {
	ok, _, err := z.client.Exists(path)
	if err != nil {
		return err
	}
	if !ok {
		if _, err := z.client.Create(path, nil, int32(flag), zk.WorldACL(zk.PermAll)); err != nil {
			return err
		}
	}
	return nil
}

func (z *ZooKeeper) generalGroupName() string {
	if z.config.GroupName == "" {
		return defaultGroupName
	}
	if strings.Index(z.config.GroupName, "/") != 0 {
		return "/" + z.config.GroupName
	}
	return z.config.GroupName
}

func (z *ZooKeeper) watcher() {
	for {
		select {
		case <-z.ctx.Done():
			return
		default:
		}

		z.targets.Range(func(key, value interface{}) bool {
			target, ok := key.(string)
			if !ok {
				return true
			}
			srvNode := fmt.Sprintf("%s/%s", z.generalGroupName(), target)
			addrs, err := z.children(srvNode)
			if err != nil {
				return true
			}
			z.targets.Store(key, addrs)

			return true
		})

		<-time.After(defaultInterval)
	}
}

func (z *ZooKeeper) children(serviceName string) ([]string, error) {
	addrs, _, err := z.client.Children(serviceName)
	if err != nil {
		return nil, err
	}
	return addrs, nil
}
