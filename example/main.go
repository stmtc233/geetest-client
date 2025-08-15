// example/main.go
package main

import (
	"context"
	"flag" // 导入 flag 包
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/stmtc233/geetest-client/geetest" // 请调整为你的模块导入路径
)

// --- 通用配置 ---
const (
	// 您要测试的目标网站的 URL。
	targetTestURL = "https://passport.bilibili.com/x/passport-login/captcha?source=main_web"
	// 您的 Rust API 服务器地址。
	apiServerAddr = "http://127.0.0.1:3000"
)

// --- 并发测试专用配置 ---
const (
	// 总共要执行的任务（请求）数量
	totalTasks = 100
	// 并发数量限制（同时运行的 goroutine 数量）
	concurrencyLimit = 50
)

// main 函数现在作为调度器，根据命令行参数选择运行模式
func main() {
	// 定义一个名为 "mode" 的字符串标志，默认值为 "single"
	// 用法信息会告诉用户可选的值
	mode := flag.String("mode", "single", "运行模式: 'single' (单个基本示例) 或 'stress' (并发压力测试)")
	flag.Parse() // 解析命令行传入的参数

	// 根据 -mode 参数的值来决定执行哪个函数
	switch *mode {
	case "single":
		runSingleExample()
	case "stress":
		runStressTest()
	default:
		log.Fatalf("无效的模式 '%s'. 请使用 'single' 或 'stress'.", *mode)
	}
}

// ===================================================================================
//
//	模式 1: 单个基本示例 (原先的代码)
//
// ===================================================================================
func runSingleExample() {
	ctx := context.Background()

	// --- 示例 1: 使用新会话的基本用法 ---
	log.Println("--- 运行模式: 单个基本示例 ---")

	// 创建一个具有特定会话 ID 的客户端
	// client, err := geetest.NewClient(apiServerAddr, geetest.WithSessionID("my-go-session-123"))
	// 创建一个带有代理的客户端
	// client, err := geetest.NewClient(apiServerAddr, geetest.WithProxy("http://user:pass@host:port"))
	// 创建一个客户端
	client, err := geetest.NewClient(apiServerAddr)
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}

	// 检查服务器是否健康
	if err := client.HealthCheck(ctx); err != nil {
		log.Fatalf("健康检查失败: %v", err)
	}
	log.Println("✅ 服务器健康。")

	// 1. 注册并获取 gt/challenge
	log.Println("1. 调用 RegisterTest...")
	regInfo, err := client.Click.RegisterTest(ctx, targetTestURL)
	if err != nil {
		log.Fatalf("RegisterTest 失败: %v", err)
	}
	gt, challenge := regInfo.First, regInfo.Second
	log.Printf("   - 获取到 GT: %s\n", gt)
	log.Printf("   - 获取到 Challenge: %s\n", challenge)

	// 2. 执行 simple match 来获取验证密钥
	log.Println("2. 调用 SimpleMatchRetry...")
	validateKey, err := client.Click.SimpleMatchRetry(ctx, gt, challenge)
	if err != nil {
		log.Fatalf("SimpleMatchRetry 失败: %v", err)
	}
	log.Printf("   - 获取到验证密钥: %s...\n", validateKey[:30])

	log.Println("--- 示例完成 ---")
}

// ===================================================================================
//
//	模式 2: 并发压力测试
//
// ===================================================================================
func runStressTest() {
	log.Printf("--- 运行模式: 并发压力测试 ---")
	log.Printf("总任务数: %d, 并发数: %d\n", totalTasks, concurrencyLimit)

	ctx := context.Background()
	var wg sync.WaitGroup
	var successCount, failureCount int32

	// 使用一个带缓冲的 channel 作为信号量，来限制并发数量
	semaphore := make(chan struct{}, concurrencyLimit)

	startTime := time.Now()

	for i := 1; i <= totalTasks; i++ {
		wg.Add(1)

		go func(taskID int) {
			defer wg.Done()

			// 获取一个信号量槽位
			semaphore <- struct{}{}
			defer func() { <-semaphore }() // 释放槽位

			log.Printf("[任务 %d] 开始执行...", taskID)
			err := runSingleTask(ctx, taskID) // 注意: 这里调用的是 runSingleTask

			if err != nil {
				log.Printf("[任务 %d] 失败: %v", taskID, err)
				atomic.AddInt32(&failureCount, 1)
			} else {
				log.Printf("[任务 %d] 成功!", taskID)
				atomic.AddInt32(&successCount, 1)
			}

		}(i)
		time.Sleep(10 * time.Microsecond)
	}

	// 等待所有 goroutine 完成
	wg.Wait()

	// --- 打印测试总结 ---
	duration := time.Since(startTime)
	totalSuccess := atomic.LoadInt32(&successCount)
	totalFailure := atomic.LoadInt32(&failureCount)
	rps := float64(totalTasks) / duration.Seconds()

	log.Println("\n--- 测试总结 ---")
	log.Printf("总耗时: %.2f 秒\n", duration.Seconds())
	log.Printf("总请求数: %d\n", totalTasks)
	log.Printf("成功请求: %d\n", totalSuccess)
	log.Printf("失败请求: %d\n", totalFailure)
	if totalTasks > 0 {
		successRate := float64(totalSuccess) / float64(totalTasks) * 100
		log.Printf("成功率: %.2f%%\n", successRate)
	}
	log.Printf("每秒请求数 (RPS): %.2f\n", rps)
	log.Println("--- 测试完成 ---")
}

// runSingleTask 封装了单次完整的点选验证流程，供并发测试使用。
// (之前名为 runSingleTest，为避免与 runStressTest 混淆，重命名为 runSingleTask)
func runSingleTask(ctx context.Context, taskID int) error {
	sessionID := fmt.Sprintf("go-stress-test-%d-%d", taskID, time.Now().UnixNano())
	// client, err := geetest.NewClient(apiServerAddr, geetest.WithSessionID(sessionID), geetest.WithProxy("http://user:pass@host:port"))
	client, err := geetest.NewClient(apiServerAddr, geetest.WithSessionID(sessionID))
	if err != nil {
		return fmt.Errorf("创建客户端失败: %w", err)
	}

	// 1. 注册
	regInfo, err := client.Click.RegisterTest(ctx, targetTestURL)
	if err != nil {
		return fmt.Errorf("RegisterTest 失败: %w", err)
	}
	gt, challenge := regInfo.First, regInfo.Second

	// 2. 获取验证结果
	_, err = client.Click.SimpleMatchRetry(ctx, gt, challenge)
	if err != nil {
		return fmt.Errorf("SimpleMatchRetry 失败: %w", err)
	}

	return nil
}
