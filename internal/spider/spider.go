package spider

import (
	"context"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Spider 语雀爬虫
type Spider struct {
	downloader       *Downloader
	progressCallback func(DownloadProgress)
	config           Config
}

// NewSpider 创建新的爬虫
func NewSpider(cookie, outputPath string, config Config, progressCallback func(DownloadProgress)) *Spider {
	return &Spider{
		downloader:       NewDownloader(cookie, outputPath, config),
		progressCallback: progressCallback,
		config:           config,
	}
}

// Download 下载知识库
func (s *Spider) Download(ctx context.Context, task DownloadTask) error {
	progress := DownloadProgress{
		Status:    "downloading",
		StartTime: time.Now(),
	}

	// 获取知识库标题
	fetcher := NewFetcher(task.Cookie, task.Config)
	bookTitle, err := fetcher.FetchBookTitle(task.URL)
	if err != nil {
		progress.Status = "error"
		progress.Error = fmt.Sprintf("获取知识库标题失败: %v", err)
		s.notifyProgress(progress)
		return err
	}

	progress.BookTitle = strings.TrimSpace(bookTitle)
	s.notifyProgress(progress)

	// 获取知识库数据
	yuqueData, err := fetcher.FetchBookData(task.URL)
	if err != nil {
		progress.Status = "error"
		progress.Error = fmt.Sprintf("获取知识库数据失败: %v", err)
		s.notifyProgress(progress)
		return err
	}

	// 解析知识库显示标题和存储目录
	displayTitle := strings.TrimSpace(yuqueData.Book.Name)
	if displayTitle == "" {
		displayTitle = progress.BookTitle
	}
	progress.BookTitle = displayTitle

	folderName := resolveBookFolderName(displayTitle, bookTitle, yuqueData.Book.ID)

	// 创建输出目录
	bookDir := filepath.Join(task.OutputPath, folderName)
	if err := os.MkdirAll(bookDir, 0755); err != nil {
		progress.Status = "error"
		progress.Error = fmt.Sprintf("创建目录失败: %v", err)
		s.notifyProgress(progress)
		return err
	}

	s.downloader.outputPath = bookDir

	// 构建目录树
	tocTree := s.buildTOCTree(yuqueData.Book.TOC)
	progress.TotalDocs = len(yuqueData.Book.TOC)
	s.notifyProgress(progress)

	// 生成 SUMMARY.md 内容
	var summaryBuilder strings.Builder

	// 下载所有文档
	for _, node := range yuqueData.Book.TOC {
		// 检查是否被取消
		select {
		case <-ctx.Done():
			progress.Status = "cancelled"
			progress.Error = "下载已取消"
			s.notifyProgress(progress)
			return ctx.Err()
		default:
		}

		progress.CurrentDoc = node.Title
		s.notifyProgress(progress)

		// 构建路径
		nodePath := tocTree[node.UUID]

		if node.Type == "TITLE" || node.ChildUUID != "" {
			// 目录节点
			if strings.HasSuffix(nodePath, "/") {
				summaryBuilder.WriteString(fmt.Sprintf("## %s\n", strings.TrimSuffix(nodePath, "/")))
			} else {
				indent := strings.Repeat("  ", strings.Count(nodePath, "/")-1)
				lastPart := nodePath[strings.LastIndex(nodePath, "/")+1:]
				summaryBuilder.WriteString(fmt.Sprintf("%s* %s\n", indent, lastPart))
			}

			// 创建目录
			dirPath := filepath.Join(bookDir, nodePath)
			os.MkdirAll(dirPath, 0755)
		}

		if node.URL != "" {
			// 文档节点
			var parentPath string
			if node.ParentUUID != "" {
				parentPath = tocTree[node.ParentUUID]
			}

			// 保存文档
			if err := s.downloader.SaveDocument(yuqueData.Book.ID, node.URL, node.Title, parentPath); err != nil {
				fmt.Printf("下载文档失败 %s: %v\n", node.Title, err)
				continue
			}

			// 添加到 SUMMARY
			indent := strings.Repeat("  ", strings.Count(parentPath, "/"))
			encodedPath := url.PathEscape(filepath.Join(parentPath, cleanFileName(node.Title)+".md"))
			summaryBuilder.WriteString(fmt.Sprintf("%s* [%s](%s)\n", indent, node.Title, encodedPath))

			progress.FinishedDocs++
			if progress.TotalDocs > 0 {
				progress.Percentage = float64(progress.FinishedDocs) / float64(progress.TotalDocs) * 100
			}
			s.notifyProgress(progress)

			// 随机延迟
			minDelay := s.config.DelayMin
			maxDelay := s.config.DelayMax
			if maxDelay > minDelay {
				delay := rand.Intn(maxDelay-minDelay) + minDelay
				time.Sleep(time.Duration(delay) * time.Second)
			}
		}
	}

	// 保存 SUMMARY.md
	summaryPath := filepath.Join(bookDir, "SUMMARY.md")
	if err := os.WriteFile(summaryPath, []byte(summaryBuilder.String()), 0644); err != nil {
		progress.Status = "error"
		progress.Error = fmt.Sprintf("保存 SUMMARY.md 失败: %v", err)
		s.notifyProgress(progress)
		return err
	}

	progress.Status = "completed"
	s.notifyProgress(progress)

	return nil
}

func resolveBookFolderName(displayTitle, fallbackTitle string, bookID int) string {
	candidates := []string{displayTitle, fallbackTitle, fmt.Sprintf("yuque-book-%d", bookID)}

	for _, candidate := range candidates {
		trimmed := strings.TrimSpace(candidate)
		if trimmed == "" {
			continue
		}

		cleaned := cleanFileName(trimmed)
		cleaned = strings.Trim(cleaned, " _-")
		if cleaned != "" {
			return cleaned
		}
	}

	return fmt.Sprintf("yuque-book-%d", bookID)
}

// buildTOCTree 构建目录树
func (s *Spider) buildTOCTree(toc []TOCNode) map[string]string {
	tree := make(map[string]string)
	nodeMap := make(map[string]*TOCNode)

	// 构建节点映射
	for i := range toc {
		node := &toc[i]
		nodeMap[node.UUID] = node
		tree[node.UUID] = ""
	}

	// 构建路径
	for _, node := range toc {
		if node.Type == "TITLE" || node.ChildUUID != "" {
			path := s.buildNodePath(node.UUID, nodeMap)
			tree[node.UUID] = path
		}
	}

	return tree
}

// buildNodePath 构建节点路径
func (s *Spider) buildNodePath(uuid string, nodeMap map[string]*TOCNode) string {
	node := nodeMap[uuid]
	if node == nil {
		return ""
	}

	parts := []string{cleanFileName(node.Title)}

	// 递归构建父路径
	current := node
	for current.ParentUUID != "" {
		parent := nodeMap[current.ParentUUID]
		if parent == nil {
			break
		}
		parts = append([]string{cleanFileName(parent.Title)}, parts...)
		current = parent
	}

	return strings.Join(parts, "/") + "/"
}

// notifyProgress 通知进度
func (s *Spider) notifyProgress(progress DownloadProgress) {
	if s.progressCallback != nil {
		s.progressCallback(progress)
	}
}
