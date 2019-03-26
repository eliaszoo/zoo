package etcd

import (
	"sync"
	"context"
	"comm/registry"
	"github.com/coreos/etcd/client"
)

type etcdRegistry struct {
	kapi   	client.KeysAPI
	nodes 	map[string]*registry.Service
	name 	string
	tag 	string

	sync.Mutex
}

func (r *etcdRegistry) Init() error {

}

func (r *etcdRegistry) Register(service *registry.Service, opts *registry.RegisterOptions) error {
	if len(s.Nodes) == 0 {
		return errors.New("Require at least one node")
	}

	var options registry.RegisterOptions
	for _, o := range opts {
		o(&options)
	}

	service := &registry.Service{
		Name:      s.Name,
		Version:   s.Version,
		Metadata:  s.Metadata,
		Endpoints: s.Endpoints,
	}

	ctx, cancel := context.WithTimeout(context.Background(), e.options.Timeout)
	defer cancel()

	_, err := e.client.Set(ctx, servicePath(s.Name), "", &etcd.SetOptions{PrevExist: etcd.PrevIgnore, Dir: true})
	if err != nil && !strings.HasPrefix(err.Error(), "102: Not a file") {
		return err
	}

	for _, node := range s.Nodes {
		service.Nodes = []*registry.Node{node}
		_, err := e.client.Set(ctx, nodePath(service.Name, node.Id), encode(service), &etcd.SetOptions{TTL: options.TTL})
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *etcdRegistry) UnRegister(service *Service) error {

}

func (r *etcdRegistry) GetService(name string) *Service {

}

func (r *etcdRegistry) ListService() []*Service {

}

func (r *etcdRegistry) Watch(opts *WatchOptions) error {
	watcher := r.kapi.Watcher(r.tag, &client.WatcherOptions{
        Recursive: true,
    })
    for {
        resp, err := watcher.Next(context.Background())
        if err != nil {
            log.Println(err)
            m.active = false
            continue
        }
        m.active = true
        switch resp.Action {
        case "set", "update":
            m.addNode(resp.Node.Key, resp.Node.Value)
            break
        case "expire", "delete":
            m.delNode(resp.Node.Key)
            break
        default:
            log.Println("watchme!!!", "resp ->", resp)
        }
    }
}