# 网页爬虫服务项目结构

```
crawler/
├── cmd/                        # 应用程序入口
│   └── server/                 # 服务器启动入口
│       └── main.go             # 主函数
├── config/                     # 配置文件
│   ├── config.go               # 配置加载与管理
│   └── config.yaml             # 配置文件
├── internal/                   # 私有应用程序代码
│   ├── api/                    # API层
│   │   ├── handlers/           # 请求处理器
│   │   │   ├── crawl.go        # 爬取任务请求处理
│   │   │   └── task.go         # 任务状态请求处理
│   │   ├── middleware/         # 中间件
│   │   │   ├── logger.go       # 日志中间件
│   │   │   └── validator.go    # 请求验证中间件
│   │   ├── models/             # API请求/响应模型
│   │   │   ├── request.go      # 请求模型
│   │   │   └── response.go     # 响应模型
│   │   └── router.go           # 路由配置
│   ├── crawler/                # 爬虫核心模块
│   │   ├── collector.go        # Colly爬虫配置
│   │   ├── domain.go           # 域名处理
│   │   └── urls.go             # URL处理与去重
│   ├── service/                # 服务层
│   │   ├── crawler_service.go  # 爬虫服务
│   │   └── task_service.go     # 任务管理服务
│   └── storage/                # 存储层
│       ├── inmemory/           # 内存存储实现
│       │   └── storage.go      # 内存存储
│       └── interface.go        # 存储接口定义
├── pkg/                        # 可重用公共代码
│   ├── logger/                 # 日志工具
│   │   └── logger.go           # 日志实现
│   └── utils/                  # 通用工具函数
│       ├── id.go               # ID生成
│       └── url.go              # URL处理
├── docs/                       # 文档
│   ├── API_SCHEMA.md           # API设计文档
│   ├── PRD.md                  # 产品需求文档
│   ├── TECH_STACK.md           # 技术栈文档
│   └── PROJECT_STRUCTURE.md    # 项目结构文档
├── test/                       # 测试
│   ├── integration/            # 集成测试
│   └── unit/                   # 单元测试
├── Dockerfile                  # Docker配置
├── go.mod                      # Go模块定义
├── go.sum                      # 依赖版本锁定
├── Makefile                    # 构建脚本
└── README.md                   # 项目说明
``` 