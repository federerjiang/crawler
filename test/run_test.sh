#!/bin/bash

# 检查参数
if [ $# -lt 1 ]; then
    echo "用法: ./run_test.sh <爬取网址>"
    echo "示例: ./run_test.sh https://example.com"
    exit 1
fi

URL=$1

# 编译并运行测试
echo "开始测试爬虫API..."
go run crawler_client.go "$URL" 