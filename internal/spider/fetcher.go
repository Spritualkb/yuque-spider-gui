package spider

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Fetcher 网络请求处理器
type Fetcher struct {
	client *http.Client
	cookie string
	config Config
}

// NewFetcher 创建新的 Fetcher
func NewFetcher(cookie string, config Config) *Fetcher {
	return &Fetcher{
		client: &http.Client{
			Timeout: time.Duration(config.Timeout) * time.Second,
		},
		cookie: cookie,
		config: config,
	}
}

// FetchBookTitle 获取知识库标题
func (f *Fetcher) FetchBookTitle(rawURL string) (string, error) {
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return "", err
	}

	if f.cookie != "" {
		req.Header.Set("Cookie", f.cookie)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("请求失败,状态码: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	title := doc.Find("title").Text()
	title = strings.TrimSpace(title)
	title = strings.ReplaceAll(title, " · 语雀", "")

	// 清理非法文件名字符
	invalidChars := regexp.MustCompile(`[\/\\:*?"<>|\n\r]`)
	title = invalidChars.ReplaceAllString(title, "-")

	// 从 URL 提取标识符
	re := regexp.MustCompile(`u\d+/([\w-]+)`)
	if matches := re.FindStringSubmatch(rawURL); len(matches) > 1 {
		title = matches[1] + "-" + title
	}

	return title, nil
}

// FetchBookData 获取知识库数据
func (f *Fetcher) FetchBookData(rawURL string) (*YuqueData, error) {
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return nil, err
	}

	if f.cookie != "" {
		req.Header.Set("Cookie", f.cookie)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 从页面中提取 JSON 数据
	re := regexp.MustCompile(`decodeURIComponent\("(.+?)"\)\);`)
	matches := re.FindStringSubmatch(string(body))
	if len(matches) < 2 {
		return nil, fmt.Errorf("无法从页面提取数据")
	}

	decodedData, err := url.QueryUnescape(matches[1])
	if err != nil {
		return nil, err
	}

	var yuqueData YuqueData
	if err := json.Unmarshal([]byte(decodedData), &yuqueData); err != nil {
		return nil, err
	}

	return &yuqueData, nil
}

// FetchDocument 获取文档内容
func (f *Fetcher) FetchDocument(bookID int, slug string) (*DocData, error) {
	apiURL := fmt.Sprintf("https://www.yuque.com/api/docs/%s?book_id=%d&merge_dynamic_data=false&mode=markdown", slug, bookID)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	if f.cookie != "" {
		req.Header.Set("Cookie", f.cookie)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("文档下载失败,状态码: %d", resp.StatusCode)
	}

	var docResp DocResponse
	if err := json.NewDecoder(resp.Body).Decode(&docResp); err != nil {
		return nil, err
	}

	return &docResp.Data, nil
}

// DownloadImage 下载图片
func (f *Fetcher) DownloadImage(imageURL string) ([]byte, error) {
	req, err := http.NewRequest("GET", imageURL, nil)
	if err != nil {
		return nil, err
	}

	if f.cookie != "" {
		req.Header.Set("Cookie", f.cookie)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("图片下载失败,状态码: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
