package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// CrawlRequest 表示爬取请求结构
type CrawlRequest struct {
	StartURL           string `json:"start_url"`
	ConcurrentRequests int    `json:"concurrent_requests,omitempty"`
	MaxDepth           int    `json:"max_depth,omitempty"`
	Timeout            int    `json:"timeout,omitempty"`
}

// Response 表示API响应的通用结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// TaskResponse 表示任务响应结构
type TaskResponse struct {
	TaskID      string    `json:"task_id"`
	Status      string    `json:"status"`
	StartURL    string    `json:"start_url"`
	Domain      string    `json:"domain,omitempty"`
	Links       []string  `json:"links,omitempty"`
	TotalLinks  int       `json:"total_links,omitempty"`
	ElapsedTime string    `json:"elapsed_time,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
}

const (
	baseURL = "http://localhost:8080/api/v1"
)

func main() {
	// 检查命令行参数
	if len(os.Args) < 2 {
		fmt.Println("使用方法: go run crawler_client.go [URL]")
		fmt.Println("示例: go run crawler_client.go https://example.com")
		os.Exit(1)
	}

	startURL := os.Args[1]
	fmt.Printf("开始爬取 URL: %s\n", startURL)

	// 创建爬取任务
	taskID, err := createCrawlTask(startURL)
	if err != nil {
		fmt.Printf("创建爬取任务失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("成功创建爬取任务，任务ID: %s\n", taskID)
	fmt.Println("等待爬取完成...")

	// 轮询检查任务状态
	for {
		task, err := getTaskStatus(taskID)
		if err != nil {
			fmt.Printf("获取任务状态失败: %v\n", err)
			break
		}

		fmt.Printf("任务状态: %s, 已爬取链接数: %d\n", task.Status, task.TotalLinks)

		if task.Status == "completed" || task.Status == "failed" {
			// 任务完成，打印爬取结果
			fmt.Printf("\n爬取结果:\n")
			fmt.Printf("起始URL: %s\n", task.StartURL)
			fmt.Printf("域名: %s\n", task.Domain)
			fmt.Printf("总链接数: %d\n", task.TotalLinks)
			fmt.Printf("耗时: %s\n", task.ElapsedTime)
			fmt.Printf("创建时间: %s\n", task.CreatedAt)
			fmt.Printf("完成时间: %s\n", task.CompletedAt)

			if len(task.Links) > 0 {
				fmt.Println("\n发现的链接:")
				for i, link := range task.Links {
					fmt.Printf("%d. %s\n", i+1, link)
				}
			}
			break
		}

		// 等待5秒后再次检查
		time.Sleep(5 * time.Second)
	}
}

// createCrawlTask 创建一个新的爬取任务
func createCrawlTask(startURL string) (string, error) {
	// 构建请求
	request := CrawlRequest{
		StartURL:           startURL,
		ConcurrentRequests: 5,  // 设置并发请求数
		MaxDepth:           3,  // 设置最大爬取深度
		Timeout:            30, // 设置超时时间（秒）
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	// 发送POST请求
	resp, err := http.Post(baseURL+"/crawl", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 解析响应
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	// 检查响应状态
	if response.Code != 200 {
		return "", fmt.Errorf("API返回错误: %s", response.Message)
	}

	// 提取任务ID
	taskData, ok := response.Data.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("无法解析任务数据")
	}

	taskID, ok := taskData["task_id"].(string)
	if !ok {
		return "", fmt.Errorf("无法获取任务ID")
	}

	return taskID, nil
}

// getTaskStatus 获取任务状态
func getTaskStatus(taskID string) (*TaskResponse, error) {
	// 发送GET请求
	resp, err := http.Get(fmt.Sprintf("%s/tasks/%s", baseURL, taskID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 解析响应
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	// 检查响应状态
	if response.Code != 200 {
		return nil, fmt.Errorf("API返回错误: %s", response.Message)
	}

	// 提取任务数据
	taskData, err := json.Marshal(response.Data)
	if err != nil {
		return nil, err
	}

	var taskResponse TaskResponse
	err = json.Unmarshal(taskData, &taskResponse)
	if err != nil {
		return nil, err
	}

	return &taskResponse, nil
}
