package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// EtcdDiscovery ETCD服务发现实现
type EtcdDiscovery struct {
	client     *clientv3.Client
	keyPrefix  string
	watchers   map[string]chan []*ServiceInfo
	mu         sync.RWMutex
	leaseID    clientv3.LeaseID
	keepAlive  <-chan *clientv3.LeaseKeepAliveResponse
	registered map[string]string // serviceID -> key
}

// NewEtcdDiscovery 创建ETCD服务发现
func NewEtcdDiscovery(endpoints []string, keyPrefix string) (*EtcdDiscovery, error) {
	if keyPrefix == "" {
		keyPrefix = "/services"
	}

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("创建ETCD客户端失败: %v", err)
	}

	return &EtcdDiscovery{
		client:     client,
		keyPrefix:  keyPrefix,
		watchers:   make(map[string]chan []*ServiceInfo),
		registered: make(map[string]string),
	}, nil
}

// Register 注册服务
func (ed *EtcdDiscovery) Register(ctx context.Context, service *ServiceInfo) error {
	// 创建租约
	lease, err := ed.client.Grant(ctx, 30) // 30秒租约
	if err != nil {
		return fmt.Errorf("创建租约失败: %v", err)
	}

	ed.leaseID = lease.ID

	// 开始保持租约活跃
	ed.keepAlive, err = ed.client.KeepAlive(ctx, ed.leaseID)
	if err != nil {
		return fmt.Errorf("保持租约失败: %v", err)
	}

	// 监听租约保活响应
	go func() {
		for ka := range ed.keepAlive {
			if ka == nil {
				break
			}
		}
	}()

	// 序列化服务信息
	data, err := json.Marshal(service)
	if err != nil {
		return fmt.Errorf("序列化服务信息失败: %v", err)
	}

	// 构建key
	key := path.Join(ed.keyPrefix, service.Name, service.ID)

	// 写入ETCD
	_, err = ed.client.Put(ctx, key, string(data), clientv3.WithLease(ed.leaseID))
	if err != nil {
		return fmt.Errorf("注册服务到ETCD失败: %v", err)
	}

	// 记录注册的服务
	ed.mu.Lock()
	ed.registered[service.ID] = key
	ed.mu.Unlock()

	return nil
}

// Deregister 注销服务
func (ed *EtcdDiscovery) Deregister(ctx context.Context, serviceID string) error {
	ed.mu.Lock()
	key, exists := ed.registered[serviceID]
	if exists {
		delete(ed.registered, serviceID)
	}
	ed.mu.Unlock()

	if !exists {
		return fmt.Errorf("服务未注册: %s", serviceID)
	}

	// 删除key
	_, err := ed.client.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("注销服务失败: %v", err)
	}

	// 撤销租约
	if ed.leaseID != 0 {
		_, err = ed.client.Revoke(ctx, ed.leaseID)
		if err != nil {
			return fmt.Errorf("撤销租约失败: %v", err)
		}
	}

	return nil
}

// Discover 发现服务
func (ed *EtcdDiscovery) Discover(ctx context.Context, serviceName string) ([]*ServiceInfo, error) {
	keyPrefix := path.Join(ed.keyPrefix, serviceName) + "/"

	resp, err := ed.client.Get(ctx, keyPrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("发现服务失败: %v", err)
	}

	var services []*ServiceInfo
	for _, kv := range resp.Kvs {
		var service ServiceInfo
		if err := json.Unmarshal(kv.Value, &service); err != nil {
			continue // 跳过无效的服务信息
		}
		services = append(services, &service)
	}

	return services, nil
}

// Watch 监听服务变化
func (ed *EtcdDiscovery) Watch(ctx context.Context, serviceName string) (<-chan []*ServiceInfo, error) {
	ed.mu.Lock()
	defer ed.mu.Unlock()

	if ch, exists := ed.watchers[serviceName]; exists {
		return ch, nil
	}

	ch := make(chan []*ServiceInfo, 10)
	ed.watchers[serviceName] = ch

	go ed.watchService(ctx, serviceName, ch)

	return ch, nil
}

// watchService 监听服务变化的goroutine
func (ed *EtcdDiscovery) watchService(ctx context.Context, serviceName string, ch chan []*ServiceInfo) {
	defer func() {
		ed.mu.Lock()
		delete(ed.watchers, serviceName)
		close(ch)
		ed.mu.Unlock()
	}()

	keyPrefix := path.Join(ed.keyPrefix, serviceName) + "/"
	watchCh := ed.client.Watch(ctx, keyPrefix, clientv3.WithPrefix())

	// 首次获取所有服务
	services, err := ed.Discover(ctx, serviceName)
	if err == nil {
		select {
		case ch <- services:
		case <-ctx.Done():
			return
		}
	}

	// 监听变化
	for {
		select {
		case <-ctx.Done():
			return
		case watchResp := <-watchCh:
			if watchResp.Err() != nil {
				continue
			}

			// 重新获取所有服务
			services, err := ed.Discover(ctx, serviceName)
			if err != nil {
				continue
			}

			select {
			case ch <- services:
			case <-ctx.Done():
				return
			}
		}
	}
}

// Close 关闭服务发现
func (ed *EtcdDiscovery) Close() error {
	ed.mu.Lock()
	defer ed.mu.Unlock()

	// 关闭所有监听器
	for _, ch := range ed.watchers {
		close(ch)
	}
	ed.watchers = make(map[string]chan []*ServiceInfo)

	// 注销所有已注册的服务
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for serviceID := range ed.registered {
		ed.Deregister(ctx, serviceID)
	}

	// 关闭ETCD客户端
	return ed.client.Close()
}
