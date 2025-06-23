package discovery

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// RoundRobinBalancer 轮询负载均衡器
type RoundRobinBalancer struct {
	services []*ServiceInfo
	current  int
	mu       sync.Mutex
}

// NewRoundRobinBalancer 创建轮询负载均衡器
func NewRoundRobinBalancer() *RoundRobinBalancer {
	return &RoundRobinBalancer{
		services: make([]*ServiceInfo, 0),
		current:  0,
	}
}

// Select 选择一个服务实例
func (rb *RoundRobinBalancer) Select(services []*ServiceInfo) (*ServiceInfo, error) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	if len(services) == 0 {
		return nil, fmt.Errorf("没有可用的服务实例")
	}

	service := services[rb.current%len(services)]
	rb.current++
	return service, nil
}

// Update 更新服务列表
func (rb *RoundRobinBalancer) Update(services []*ServiceInfo) {
	rb.mu.Lock()
	defer rb.mu.Unlock()
	rb.services = services
	rb.current = 0
}

// RandomBalancer 随机负载均衡器
type RandomBalancer struct {
	rand *rand.Rand
	mu   sync.Mutex
}

// NewRandomBalancer 创建随机负载均衡器
func NewRandomBalancer() *RandomBalancer {
	return &RandomBalancer{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Select 选择一个服务实例
func (rb *RandomBalancer) Select(services []*ServiceInfo) (*ServiceInfo, error) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	if len(services) == 0 {
		return nil, fmt.Errorf("没有可用的服务实例")
	}

	index := rb.rand.Intn(len(services))
	return services[index], nil
}

// Update 更新服务列表
func (rb *RandomBalancer) Update(services []*ServiceInfo) {
	// 随机负载均衡器不需要维护状态
}

// WeightedRoundRobinBalancer 加权轮询负载均衡器
type WeightedRoundRobinBalancer struct {
	services       []*ServiceInfo
	weights        []int
	currentWeights []int
	totalWeight    int
	mu             sync.Mutex
}

// NewWeightedRoundRobinBalancer 创建加权轮询负载均衡器
func NewWeightedRoundRobinBalancer() *WeightedRoundRobinBalancer {
	return &WeightedRoundRobinBalancer{
		services:       make([]*ServiceInfo, 0),
		weights:        make([]int, 0),
		currentWeights: make([]int, 0),
	}
}

// Select 选择一个服务实例
func (wrb *WeightedRoundRobinBalancer) Select(services []*ServiceInfo) (*ServiceInfo, error) {
	wrb.mu.Lock()
	defer wrb.mu.Unlock()

	if len(services) == 0 {
		return nil, fmt.Errorf("没有可用的服务实例")
	}

	// 如果服务列表发生变化，重新初始化权重
	if len(services) != len(wrb.services) {
		wrb.initWeights(services)
	}

	// 找到当前权重最大的服务
	maxWeight := -1
	selectedIndex := -1
	for i, weight := range wrb.currentWeights {
		if weight > maxWeight {
			maxWeight = weight
			selectedIndex = i
		}
	}

	if selectedIndex == -1 {
		return services[0], nil
	}

	// 更新当前权重
	wrb.currentWeights[selectedIndex] -= wrb.totalWeight
	for i := range wrb.currentWeights {
		wrb.currentWeights[i] += wrb.weights[i]
	}

	return services[selectedIndex], nil
}

// Update 更新服务列表
func (wrb *WeightedRoundRobinBalancer) Update(services []*ServiceInfo) {
	wrb.mu.Lock()
	defer wrb.mu.Unlock()
	wrb.initWeights(services)
}

// initWeights 初始化权重
func (wrb *WeightedRoundRobinBalancer) initWeights(services []*ServiceInfo) {
	wrb.services = services
	wrb.weights = make([]int, len(services))
	wrb.currentWeights = make([]int, len(services))
	wrb.totalWeight = 0

	for i, service := range services {
		weight := 1 // 默认权重为1
		if weightStr, exists := service.Metadata["weight"]; exists {
			if w, err := fmt.Sscanf(weightStr, "%d", &weight); err != nil || w != 1 {
				weight = 1
			}
		}
		wrb.weights[i] = weight
		wrb.currentWeights[i] = weight
		wrb.totalWeight += weight
	}
}

// ConsistentHashBalancer 一致性哈希负载均衡器
type ConsistentHashBalancer struct {
	services []*ServiceInfo
	mu       sync.RWMutex
}

// NewConsistentHashBalancer 创建一致性哈希负载均衡器
func NewConsistentHashBalancer() *ConsistentHashBalancer {
	return &ConsistentHashBalancer{
		services: make([]*ServiceInfo, 0),
	}
}

// SelectByKey 根据key选择服务实例
func (chb *ConsistentHashBalancer) SelectByKey(services []*ServiceInfo, key string) (*ServiceInfo, error) {
	chb.mu.RLock()
	defer chb.mu.RUnlock()

	if len(services) == 0 {
		return nil, fmt.Errorf("没有可用的服务实例")
	}

	// 简单的哈希算法
	hash := 0
	for _, b := range []byte(key) {
		hash = hash*31 + int(b)
	}
	if hash < 0 {
		hash = -hash
	}

	index := hash % len(services)
	return services[index], nil
}

// Select 选择一个服务实例（随机选择）
func (chb *ConsistentHashBalancer) Select(services []*ServiceInfo) (*ServiceInfo, error) {
	return chb.SelectByKey(services, fmt.Sprintf("%d", time.Now().UnixNano()))
}

// Update 更新服务列表
func (chb *ConsistentHashBalancer) Update(services []*ServiceInfo) {
	chb.mu.Lock()
	defer chb.mu.Unlock()
	chb.services = services
}
