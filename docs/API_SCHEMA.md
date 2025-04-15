# 网页爬虫服务 API Schema 设计

## API 概述

本文档定义了网页爬虫服务的RESTful API接口，包括请求/响应格式、状态码和错误处理方式。

## 基础信息

- **基础URL**: `/api/v1`
- **数据格式**: JSON
- **认证方式**: 无（可根据需求扩展）

## API 端点

### 1. 创建爬取任务

创建一个新的网页爬取任务。

- **URL**: `/crawl`
- **方法**: `POST`
- **请求体**:

```json
{
  "start_url": "https://example.com",          // 必填，起始URL
  "concurrent_requests": 10,                    // 可选，并发请求数，默认值10
  "max_depth": 5,                              // 可选，最大爬取深度，默认无限制
  "timeout": 30                                // 可选，超时时间(秒)，默认30秒
}
```

- **响应**:

```json
{
  "code": 200,                                  // 状态码
  "message": "success",                         // 状态信息
  "data": {
    "task_id": "550e8400-e29b-41d4-a716-446655440000", // 任务ID
    "status": "pending",                        // 任务状态: pending, in_progress
    "start_url": "https://example.com",         // 起始URL
    "created_at": "2023-04-20T10:30:00Z"        // 创建时间
  }
}
```

- **状态码**:
  - `200 OK`: 请求成功
  - `400 Bad Request`: 请求参数错误
  - `500 Internal Server Error`: 服务器内部错误

### 2. 获取爬取任务状态

获取指定爬取任务的状态和结果。

- **URL**: `/tasks/{taskId}`
- **方法**: `GET`
- **URL参数**:
  - `taskId`: 任务ID

- **响应**:

```json
{
  "code": 200,                                  // 状态码
  "message": "success",                         // 状态信息
  "data": {
    "task_id": "550e8400-e29b-41d4-a716-446655440000", // 任务ID
    "status": "completed",                      // 任务状态: pending, in_progress, completed, failed
    "start_url": "https://example.com",         // 起始URL
    "domain": "example.com",                    // 域名
    "links": [                                  // 爬取到的链接列表
      "https://example.com",
      "https://example.com/page1",
      "https://example.com/page2"
    ],
    "total_links": 3,                           // 链接总数
    "elapsed_time": "10.5s",                    // 耗时
    "created_at": "2023-04-20T10:30:00Z",       // 创建时间
    "completed_at": "2023-04-20T10:30:10Z"      // 完成时间
  }
}
```

- **状态码**:
  - `200 OK`: 请求成功
  - `404 Not Found`: 任务不存在
  - `500 Internal Server Error`: 服务器内部错误

## 错误响应

当API调用失败时，将返回统一格式的错误响应：

```json
{
  "code": 400,                 // HTTP状态码
  "message": "错误消息",        // 简短错误描述
  "error": {
    "reason": "详细错误原因",   // 详细错误信息
    "field": "错误字段名"       // 可选，指出哪个字段有问题
  }
}
```

## 错误代码与消息

| 代码 | 消息 | 描述 |
|-----|------|-----|
| 400 | Bad Request | 请求参数错误 |
| 404 | Not Found | 资源不存在 |
| 408 | Request Timeout | 请求超时 |
| 429 | Too Many Requests | 请求频率超限 |
| 500 | Internal Server Error | 服务器内部错误 |

## 请求参数校验规则

### 起始URL (start_url)
- 必须是有效的URL格式
- 必须包含协议（http/https）

### 并发请求数 (concurrent_requests)
- 正整数
- 范围: 1-50
- 默认值: 10

### 最大爬取深度 (max_depth)
- 正整数
- 范围: 1-100
- 默认值: 10 (无限制用-1表示)

### 超时时间 (timeout)
- 正整数
- 范围: 5-120
- 默认值: 30

## 示例调用

### 创建爬取任务

**请求**:
```bash
curl -X POST http://localhost:8080/api/v1/crawl \
  -H "Content-Type: application/json" \
  -d '{
    "start_url": "https://example.com",
    "concurrent_requests": 5,
    "max_depth": 3
  }'
```

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "task_id": "550e8400-e29b-41d4-a716-446655440000",
    "status": "pending",
    "start_url": "https://example.com",
    "created_at": "2023-04-20T10:30:00Z"
  }
}
```

### 获取任务状态

**请求**:
```bash
curl -X GET http://localhost:8080/api/v1/tasks/550e8400-e29b-41d4-a716-446655440000
```

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "task_id": "550e8400-e29b-41d4-a716-446655440000",
    "status": "completed",
    "start_url": "https://example.com",
    "domain": "example.com",
    "links": [
      "https://example.com",
      "https://example.com/page1",
      "https://example.com/page2"
    ],
    "total_links": 3,
    "elapsed_time": "10.5s",
    "created_at": "2023-04-20T10:30:00Z",
    "completed_at": "2023-04-20T10:30:10Z"
  }
}
``` 