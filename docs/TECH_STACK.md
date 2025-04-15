# 网页爬虫服务技术栈文档

## 核心技术选型

### 后端开发
| 技术 | 版本 | 用途 |
|-----|------|-----|
| Go | 1.20+ | 核心编程语言 |
| Gin | v1.9.0+ | RESTful API框架 |
| Colly | v2.1.0+ | 网页爬虫框架 |

### 工具与中间件
| 技术 | 版本 | 用途 |
|-----|------|-----|
| Docker | 20.10+ | 容器化部署 |
| Swagger | 2.0 | API文档 |

## 应用架构

### 系统组件
```
┌───────────────┐      ┌─────────────────┐      ┌─────────────────┐
│  API 层       │      │  服务层         │      │  数据存储层     │
│  (Gin)        │─────▶│  (爬虫服务)     │─────▶│  (内存)         │
└───────────────┘      └─────────────────┘      └─────────────────┘
                               │
                               ▼
                       ┌─────────────────┐
                       │  外部网站       │
                       │  (爬取目标)     │
                       └─────────────────┘
```

### 模块划分
1. **API模块**
   - 处理HTTP请求响应
   - 参数验证
   - 返回格式化结果

2. **任务管理模块**
   - 任务创建与调度
   - 状态跟踪
   - 结果存储与查询

3. **爬虫核心模块**
   - 基于Colly的爬虫实现
   - 域名过滤
   - URL去重
   - 并发控制

4. **存储模块**
   - 内存缓存实现

## 详细技术方案

### Go + Gin实现RESTful API
```go
// 主要路由设置
func setupRouter() *gin.Engine {
    router := gin.Default()
    
    v1 := router.Group("/api/v1")
    {
        v1.POST("/crawl", handlers.StartCrawl)
        v1.GET("/tasks/:taskId", handlers.GetTaskStatus)
    }
    
    return router
}
```

### Colly爬虫实现
```go
// 爬虫核心实现
func createCrawler(startURL string, concurrentRequests int) *colly.Collector {
    // 提取域名
    domain := extractDomain(startURL)
    
    // 创建爬虫实例
    c := colly.NewCollector(
        colly.AllowedDomains(domain),         // 域名限制
        colly.Async(true),                     // 启用异步
        colly.MaxDepth(10),                    // 最大深度限制
    )
    
    // 设置并发限制
    c.Limit(&colly.LimitRule{
        DomainGlob:  "*",
        Parallelism: concurrentRequests,
    })
    
    return c
}
```

### 域名过滤
```go
// 域名提取与验证
func extractDomain(url string) string {
    parsed, err := neturl.Parse(url)
    if err != nil {
        return ""
    }
    return parsed.Hostname()
}

func isSameDomain(url1, url2 string) bool {
    domain1 := extractDomain(url1)
    domain2 := extractDomain(url2)
    return domain1 == domain2
}
```

### URL去重处理
```go
// URL去重
type VisitedURLs struct {
    urls   map[string]bool
    mutex  sync.RWMutex
}

func NewVisitedURLs() *VisitedURLs {
    return &VisitedURLs{
        urls: make(map[string]bool),
    }
}

func (v *VisitedURLs) Add(url string) bool {
    v.mutex.Lock()
    defer v.mutex.Unlock()
    
    if _, exists := v.urls[url]; exists {
        return false
    }
    
    v.urls[url] = true
    return true
}
```

### 任务管理
```go
// 任务数据结构
type CrawlTask struct {
    ID               string    `json:"task_id"`
    StartURL         string    `json:"start_url"`
    Status           string    `json:"status"` // pending, in_progress, completed, failed
    Links            []string  `json:"links,omitempty"`
    TotalLinks       int       `json:"total_links"`
    ElapsedTime      string    `json:"elapsed_time,omitempty"`
    CreatedAt        time.Time `json:"created_at"`
    CompletedAt      time.Time `json:"completed_at,omitempty"`
}
```

### 存储实现
```go
// 内存存储实现
type InMemoryStorage struct {
    tasks  map[string]*CrawlTask
    mutex  sync.RWMutex
}

func NewInMemoryStorage() *InMemoryStorage {
    return &InMemoryStorage{
        tasks: make(map[string]*CrawlTask),
    }
}

func (s *InMemoryStorage) SaveTask(task *CrawlTask) error {
    s.mutex.Lock()
    defer s.mutex.Unlock()
    s.tasks[task.ID] = task
    return nil
}

func (s *InMemoryStorage) GetTask(taskID string) (*CrawlTask, error) {
    s.mutex.RLock()
    defer s.mutex.RUnlock()
    
    task, exists := s.tasks[taskID]
    if !exists {
        return nil, errors.New("task not found")
    }
    return task, nil
}

func (s *InMemoryStorage) UpdateTaskStatus(taskID, status string) error {
    s.mutex.Lock()
    defer s.mutex.Unlock()
    
    task, exists := s.tasks[taskID]
    if !exists {
        return errors.New("task not found")
    }
    
    task.Status = status
    if status == "completed" || status == "failed" {
        task.CompletedAt = time.Now()
    }
    
    return nil
}
```

## 部署方案

### Docker部署
```dockerfile
FROM golang:1.20-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o crawler ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/crawler .
COPY --from=builder /app/config.yaml .

EXPOSE 8080
CMD ["./crawler"]
```

### 配置管理
```yaml
# config.yaml
server:
  port: 8080
  timeout: 30s

crawler:
  default_concurrent_requests: 10
  default_timeout: 30s
  default_max_depth: 5
```

## 性能优化策略

1. **爬虫性能优化**
   - 使用Colly的并发功能
   - 实现自适应的爬取速率控制
   - 超时控制与重试机制

2. **API性能优化**
   - Gin框架的高效路由
   - 异步任务处理
   - 请求限流

3. **资源使用优化**
   - 内存池与对象复用
   - 定期清理过期任务
   - 限制单个任务的最大URL数
   - 实现最大任务数量限制，避免内存溢出

## 日志

1. **API日志**
   - 请求路径记录
   - 错误日志

2. **爬虫日志**
   - 爬取URL数量统计
   - 失败URL记录

## 测试策略

1. **单元测试**
   - 域名提取和验证
   - URL去重逻辑
   - 任务管理

2. **集成测试**
   - API端点测试
   - 爬虫功能测试

3. **性能测试**
   - 并发爬取测试
   - API负载测试 