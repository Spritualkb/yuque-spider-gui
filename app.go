package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"yuque-spider-gui/internal/spider"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// TaskStatus 任务状态
type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusRunning    TaskStatus = "running"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusFailed     TaskStatus = "failed"
	TaskStatusCancelled  TaskStatus = "cancelled"
)

// DownloadTaskItem 下载任务项
type DownloadTaskItem struct {
	ID          string                   `json:"id"`
	URL         string                   `json:"url"`
	Cookie      string                   `json:"cookie"`
	OutputPath  string                   `json:"outputPath"`
	Config      spider.Config            `json:"config"`
	Status      TaskStatus               `json:"status"`
	Progress    spider.DownloadProgress  `json:"progress"`
	Error       string                   `json:"error,omitempty"`
	CreatedAt   time.Time                `json:"createdAt"`
	StartedAt   *time.Time               `json:"startedAt,omitempty"`
	CompletedAt *time.Time               `json:"completedAt,omitempty"`
	cancelFunc  context.CancelFunc
	spider      *spider.Spider
}

// App struct
type App struct {
	ctx       context.Context
	tasks     map[string]*DownloadTaskItem
	taskOrder []string // 保持任务顺序
	mu        sync.RWMutex
	taskIDCounter int
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		tasks:     make(map[string]*DownloadTaskItem),
		taskOrder: make([]string, 0),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// GetDefaultConfig 获取默认配置
func (a *App) GetDefaultConfig() spider.Config {
	return spider.DefaultConfig()
}

// SelectDirectory 选择目录
func (a *App) SelectDirectory() (string, error) {
	dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择下载目录",
	})
	return dir, err
}

// AddTask 添加任务
func (a *App) AddTask(url, cookie, outputPath string, config spider.Config) (string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	// 生成任务 ID
	a.taskIDCounter++
	taskID := fmt.Sprintf("task_%d", a.taskIDCounter)

	// 验证输出路径
	if outputPath == "" {
		homeDir, _ := os.UserHomeDir()
		outputPath = filepath.Join(homeDir, "Downloads", "yuque-downloads")
	}

	// 设置默认配置
	if config.Timeout == 0 {
		config = spider.DefaultConfig()
	}

	// 创建任务
	task := &DownloadTaskItem{
		ID:         taskID,
		URL:        url,
		Cookie:     cookie,
		OutputPath: outputPath,
		Config:     config,
		Status:     TaskStatusPending,
		CreatedAt:  time.Now(),
		Progress: spider.DownloadProgress{
			Status: "pending",
		},
	}

	a.tasks[taskID] = task
	a.taskOrder = append(a.taskOrder, taskID)

	// 通知前端任务列表更新
	a.emitTaskListUpdate()

	return taskID, nil
}

// RemoveTask 删除任务
func (a *App) RemoveTask(taskID string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	task, exists := a.tasks[taskID]
	if !exists {
		return fmt.Errorf("任务不存在: %s", taskID)
	}

	// 如果任务正在运行，先取消
	if task.Status == TaskStatusRunning && task.cancelFunc != nil {
		task.cancelFunc()
	}

	// 删除任务
	delete(a.tasks, taskID)

	// 从顺序列表中删除
	for i, id := range a.taskOrder {
		if id == taskID {
			a.taskOrder = append(a.taskOrder[:i], a.taskOrder[i+1:]...)
			break
		}
	}

	a.emitTaskListUpdate()
	return nil
}

// StartTask 开始任务
func (a *App) StartTask(taskID string) error {
	a.mu.Lock()
	task, exists := a.tasks[taskID]
	if !exists {
		a.mu.Unlock()
		return fmt.Errorf("任务不存在: %s", taskID)
	}

	if task.Status == TaskStatusRunning {
		a.mu.Unlock()
		return fmt.Errorf("任务正在运行中")
	}

	// 确保输出目录存在
	if err := os.MkdirAll(task.OutputPath, 0755); err != nil {
		a.mu.Unlock()
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	// 更新任务状态
	task.Status = TaskStatusRunning
	now := time.Now()
	task.StartedAt = &now
	task.Error = ""

	// 创建上下文
	ctx, cancel := context.WithCancel(a.ctx)
	task.cancelFunc = cancel

	// 创建爬虫实例
	task.spider = spider.NewSpider(task.Cookie, task.OutputPath, task.Config, func(progress spider.DownloadProgress) {
		a.mu.Lock()
		defer a.mu.Unlock()

		if t, ok := a.tasks[taskID]; ok {
			t.Progress = progress

			// 检查是否完成
			if progress.Status == "completed" {
				t.Status = TaskStatusCompleted
				now := time.Now()
				t.CompletedAt = &now
			} else if progress.Status == "error" || progress.Status == "cancelled" {
				if progress.Status == "error" {
					t.Status = TaskStatusFailed
					t.Error = progress.Error
				} else {
					t.Status = TaskStatusCancelled
				}
				now := time.Now()
				t.CompletedAt = &now
			}

			// 发送任务更新事件
			runtime.EventsEmit(a.ctx, "task:update", t)
		}
	})

	a.mu.Unlock()

	// 在后台启动下载
	go func() {
		downloadTask := spider.DownloadTask{
			URL:        task.URL,
			Cookie:     task.Cookie,
			OutputPath: task.OutputPath,
			Config:     task.Config,
		}

		err := task.spider.Download(ctx, downloadTask)

		a.mu.Lock()
		defer a.mu.Unlock()

		if t, ok := a.tasks[taskID]; ok {
			if err != nil && t.Status == TaskStatusRunning {
				t.Status = TaskStatusFailed
				t.Error = err.Error()
				now := time.Now()
				t.CompletedAt = &now
				runtime.EventsEmit(a.ctx, "task:update", t)
			}
		}
	}()

	return nil
}

// CancelTask 取消任务
func (a *App) CancelTask(taskID string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	task, exists := a.tasks[taskID]
	if !exists {
		return fmt.Errorf("任务不存在: %s", taskID)
	}

	if task.Status != TaskStatusRunning {
		return fmt.Errorf("任务未在运行中")
	}

	if task.cancelFunc != nil {
		task.cancelFunc()
		task.Status = TaskStatusCancelled
		now := time.Now()
		task.CompletedAt = &now
		runtime.EventsEmit(a.ctx, "task:update", task)
	}

	return nil
}

// GetAllTasks 获取所有任务
func (a *App) GetAllTasks() []DownloadTaskItem {
	a.mu.RLock()
	defer a.mu.RUnlock()

	tasks := make([]DownloadTaskItem, 0, len(a.taskOrder))
	for _, taskID := range a.taskOrder {
		if task, exists := a.tasks[taskID]; exists {
			// 创建副本，不暴露内部字段
			taskCopy := *task
			taskCopy.cancelFunc = nil
			taskCopy.spider = nil
			tasks = append(tasks, taskCopy)
		}
	}

	return tasks
}

// ClearCompletedTasks 清除已完成的任务
func (a *App) ClearCompletedTasks() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	newTaskOrder := make([]string, 0)
	for _, taskID := range a.taskOrder {
		task := a.tasks[taskID]
		if task.Status == TaskStatusCompleted || task.Status == TaskStatusFailed || task.Status == TaskStatusCancelled {
			delete(a.tasks, taskID)
		} else {
			newTaskOrder = append(newTaskOrder, taskID)
		}
	}

	a.taskOrder = newTaskOrder
	a.emitTaskListUpdate()
	return nil
}

// StartAllPendingTasks 开始所有待处理任务
func (a *App) StartAllPendingTasks() error {
	a.mu.RLock()
	pendingTasks := make([]string, 0)
	for _, taskID := range a.taskOrder {
		if task := a.tasks[taskID]; task.Status == TaskStatusPending {
			pendingTasks = append(pendingTasks, taskID)
		}
	}
	a.mu.RUnlock()

	for _, taskID := range pendingTasks {
		if err := a.StartTask(taskID); err != nil {
			return err
		}
		// 稍微延迟，避免同时启动太多任务
		time.Sleep(500 * time.Millisecond)
	}

	return nil
}

// ValidateURL 验证 URL 是否有效
func (a *App) ValidateURL(url string) bool {
	return len(url) > 0 && (contains(url, "yuque.com") || contains(url, "www.yuque.com"))
}

// emitTaskListUpdate 发送任务列表更新事件
func (a *App) emitTaskListUpdate() {
	tasks := make([]DownloadTaskItem, 0, len(a.taskOrder))
	for _, taskID := range a.taskOrder {
		if task, exists := a.tasks[taskID]; exists {
			taskCopy := *task
			taskCopy.cancelFunc = nil
			taskCopy.spider = nil
			tasks = append(tasks, taskCopy)
		}
	}
	runtime.EventsEmit(a.ctx, "tasks:update", tasks)
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsSubstring(s, substr)))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
