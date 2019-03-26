package registry

type Registry interface {
	Init(opts *Options)
	Register(service *Service, opts *RegisterOptions) error
	UnRegister(service *Service) error
	GetService(name string) *Service
	ListService() []*Service
	Watch(opts *WatchOptions) error
}