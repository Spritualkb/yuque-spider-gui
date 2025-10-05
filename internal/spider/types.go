package spider

import "time"

// Config 爬虫配置
type Config struct {
	// DelayMin 最小延迟(秒)
	DelayMin int `json:"delayMin"`
	// DelayMax 最大延迟(秒)
	DelayMax int `json:"delayMax"`
	// Timeout 请求超时时间(秒)
	Timeout int `json:"timeout"`
	// MaxRetries 最大重试次数
	MaxRetries int `json:"maxRetries"`
	// ConcurrentDownloads 并发下载数
	ConcurrentDownloads int `json:"concurrentDownloads"`
}

// DefaultConfig 默认配置
func DefaultConfig() Config {
	return Config{
		DelayMin:            1,
		DelayMax:            4,
		Timeout:             30,
		MaxRetries:          3,
		ConcurrentDownloads: 1,
	}
}

// Book 语雀知识库结构
type Book struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	TOC         []TOCNode `json:"toc"`
}

// TOCNode 目录节点
type TOCNode struct {
	UUID       string `json:"uuid"`
	Title      string `json:"title"`
	URL        string `json:"url"`
	Slug       string `json:"slug"`
	Type       string `json:"type"`
	ParentUUID string `json:"parent_uuid"`
	ChildUUID  string `json:"child_uuid"`
	Depth      int    `json:"depth"`
}

// DocResponse 文档 API 响应
type DocResponse struct {
	Data DocData `json:"data"`
}

// DocData 文档数据
type DocData struct {
	ID         int    `json:"id"`
	Slug       string `json:"slug"`
	Title      string `json:"title"`
	SourceCode string `json:"sourcecode"`
}

// YuqueData 页面中的数据
type YuqueData struct {
	Book Book `json:"book"`
}

// DownloadTask 下载任务
type DownloadTask struct {
	URL        string `json:"url"`
	Cookie     string `json:"cookie"`
	OutputPath string `json:"outputPath"`
	Config     Config `json:"config"`
}

// DownloadProgress 下载进度
type DownloadProgress struct {
	BookTitle    string    `json:"bookTitle"`
	CurrentDoc   string    `json:"currentDoc"`
	TotalDocs    int       `json:"totalDocs"`
	FinishedDocs int       `json:"finishedDocs"`
	Status       string    `json:"status"` // downloading, completed, error, cancelled
	Error        string    `json:"error,omitempty"`
	StartTime    time.Time `json:"startTime"`
	Percentage   float64   `json:"percentage"`
}
