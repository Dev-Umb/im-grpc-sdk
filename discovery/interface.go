package discovery

import (
	"context"
	"time"
)

// ServiceInfo 服务信息
type ServiceInfo struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Address  string            `json:"address"`
	Port     int               `json:"port"`
	Tags     []string          `json:"tags"`
	Metadata map[string]string `json:"metadata"`
	Health   string            `json:"health"` // healthy, unhealthy, unknown
}

// ServiceDiscovery 服务发现接口
type ServiceDiscovery interface {
	// Register 注册服务
	Register(ctx context.Context, service *ServiceInfo) error

	// Deregister 注销服务
	Deregister(ctx context.Context, serviceID string) error

	// Discover 发现服务
	Discover(ctx context.Context, serviceName string) ([]*ServiceInfo, error)

	// Watch 监听服务变化
	Watch(ctx context.Context, serviceName string) (<-chan []*ServiceInfo, error)

	// Close 关闭服务发现
	Close() error
}

// LoadBalancer 负载均衡器接口
type LoadBalancer interface {
	// Select 选择一个服务实例
	Select(services []*ServiceInfo) (*ServiceInfo, error)

	// Update 更新服务列表
	Update(services []*ServiceInfo)
}

// HealthChecker 健康检查接口
type HealthChecker interface {
	// Check 检查服务健康状态
	Check(ctx context.Context, service *ServiceInfo) error

	// StartHealthCheck 开始健康检查
	StartHealthCheck(ctx context.Context, service *ServiceInfo, interval time.Duration) <-chan error

	// StopHealthCheck 停止健康检查
	StopHealthCheck(serviceID string)
}
