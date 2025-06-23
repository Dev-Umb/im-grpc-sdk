package discovery

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/consul/api"
)

// ConsulDiscovery Consul服务发现实现
type ConsulDiscovery struct {
	client   *api.Client
	config   *api.Config
	watchers map[string]chan []*ServiceInfo
	mu       sync.RWMutex
}

// NewConsulDiscovery 创建Consul服务发现
func NewConsulDiscovery(address string) (*ConsulDiscovery, error) {
	config := api.DefaultConfig()
	if address != "" {
		config.Address = address
	}

	client, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("创建Consul客户端失败: %v", err)
	}

	return &ConsulDiscovery{
		client:   client,
		config:   config,
		watchers: make(map[string]chan []*ServiceInfo),
	}, nil
}

// Register 注册服务
func (cd *ConsulDiscovery) Register(ctx context.Context, service *ServiceInfo) error {
	registration := &api.AgentServiceRegistration{
		ID:      service.ID,
		Name:    service.Name,
		Address: service.Address,
		Port:    service.Port,
		Tags:    service.Tags,
		Meta:    service.Metadata,
		Check: &api.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%s:%d", service.Address, service.Port),
			Interval:                       "10s",
			Timeout:                        "5s",
			DeregisterCriticalServiceAfter: "30s",
		},
	}

	return cd.client.Agent().ServiceRegister(registration)
}

// Deregister 注销服务
func (cd *ConsulDiscovery) Deregister(ctx context.Context, serviceID string) error {
	return cd.client.Agent().ServiceDeregister(serviceID)
}

// Discover 发现服务
func (cd *ConsulDiscovery) Discover(ctx context.Context, serviceName string) ([]*ServiceInfo, error) {
	services, _, err := cd.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, fmt.Errorf("发现服务失败: %v", err)
	}

	var result []*ServiceInfo
	for _, service := range services {
		info := &ServiceInfo{
			ID:       service.Service.ID,
			Name:     service.Service.Service,
			Address:  service.Service.Address,
			Port:     service.Service.Port,
			Tags:     service.Service.Tags,
			Metadata: service.Service.Meta,
			Health:   "healthy", // 只返回健康的服务
		}
		result = append(result, info)
	}

	return result, nil
}

// Watch 监听服务变化
func (cd *ConsulDiscovery) Watch(ctx context.Context, serviceName string) (<-chan []*ServiceInfo, error) {
	cd.mu.Lock()
	defer cd.mu.Unlock()

	if ch, exists := cd.watchers[serviceName]; exists {
		return ch, nil
	}

	ch := make(chan []*ServiceInfo, 10)
	cd.watchers[serviceName] = ch

	go cd.watchService(ctx, serviceName, ch)

	return ch, nil
}

// watchService 监听服务变化的goroutine
func (cd *ConsulDiscovery) watchService(ctx context.Context, serviceName string, ch chan []*ServiceInfo) {
	defer func() {
		cd.mu.Lock()
		delete(cd.watchers, serviceName)
		close(ch)
		cd.mu.Unlock()
	}()

	var lastIndex uint64
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			services, meta, err := cd.client.Health().Service(serviceName, "", true, &api.QueryOptions{
				WaitIndex: lastIndex,
				WaitTime:  30 * time.Second,
			})
			if err != nil {
				continue
			}

			if meta.LastIndex > lastIndex {
				lastIndex = meta.LastIndex
				var result []*ServiceInfo
				for _, service := range services {
					info := &ServiceInfo{
						ID:       service.Service.ID,
						Name:     service.Service.Service,
						Address:  service.Service.Address,
						Port:     service.Service.Port,
						Tags:     service.Service.Tags,
						Metadata: service.Service.Meta,
						Health:   "healthy",
					}
					result = append(result, info)
				}

				select {
				case ch <- result:
				case <-ctx.Done():
					return
				}
			}
		}
	}
}

// Close 关闭服务发现
func (cd *ConsulDiscovery) Close() error {
	cd.mu.Lock()
	defer cd.mu.Unlock()

	for _, ch := range cd.watchers {
		close(ch)
	}
	cd.watchers = make(map[string]chan []*ServiceInfo)

	return nil
}
