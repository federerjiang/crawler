# 网页爬虫服务需求文档

## 项目概述
开发一个基于Go语言和Colly框架的RESTful API服务，该服务能够根据用户提供的起始URL，自动爬取同一域名下的所有网页，并返回完整的网页链接列表。

## 功能需求

### 核心功能
1. 接收用户提供的起始URL（startURL）
2. 自动爬取startURL同一域名下的所有网页
3. 返回所有爬取到的网页链接列表
4. 过滤掉非同一域名的外部链接
5. 处理重复链接（去重）
6. 支持并发爬取以提高效率

### API规范
1. RESTful API设计
   - `POST /api/v1/crawl` - 开始爬取任务
   - `GET /api/v1/tasks/{taskId}` - 获取爬取任务状态及结果

2. 请求/响应格式
   ```json
   // 爬取请求
   POST /api/v1/crawl
   {
     "start_url": "https://example.com",
     "concurrent_requests": 10, // 可选，默认值10
     "max_depth": 5,           // 可选，默认值无限制
     "timeout": 30             // 可选，单位秒，默认30秒
   }

   // 响应
   {
     "task_id": "12345",
     "status": "in_progress"
   }

   // 获取结果
   GET /api/v1/tasks/12345

   // 响应
   {
     "task_id": "12345",
     "status": "completed",
     "start_url": "https://example.com",
     "domain": "example.com",
     "links": [
       "https://example.com",
       "https://example.com/page1",
       "https://example.com/page2",
       ...
     ],
     "total_links": 42,
     "elapsed_time": "10.5s",
     "created_at": "2023-04-20T10:30:00Z",
     "completed_at": "2023-04-20T10:30:10Z"
   }
   ```

## 技术需求

### 技术栈
- 编程语言: Go
- 爬虫框架: Colly
- API框架: Gin
- 存储: 内存缓存或Redis (任务状态和结果)

### 实现细节
1. **域名限制**:
   - 提取startURL的域名
   - 只爬取与startURL相同域名的页面
   - 过滤掉外部链接

2. **去重处理**:
   - 使用集合或哈希表存储已访问的URL
   - 避免重复爬取同一网页

3. **并发控制**:
   - 利用Colly的并发功能
   - 可配置最大并发请求数
   - 使用适当的超时和错误处理

4. **异步任务处理**:
   - 爬取任务在后台异步执行
   - 提供任务ID用于查询进度和结果

## 性能要求
1. 服务应能处理中等规模网站的爬取（数千个网页）
2. API响应时间应在合理范围内（创建任务响应<1秒）
3. 合理控制爬取速率，避免对目标网站造成过大压力
4. 内存使用应受控，避免大型爬取任务耗尽系统资源

## 部署要求
1. 提供Docker容器化支持
2. 配置文件支持环境变量注入
3. 详细的API文档（Swagger/OpenAPI）
4. 基本的监控和日志记录

## 后续优化方向
1. 爬取结果持久化存储
2. 增加更多过滤选项（如robots.txt遵循、URL模式匹配等）
3. 支持定时爬取和变更检测
4. 增加爬取内容解析和数据提取功能
