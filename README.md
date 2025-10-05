# 语雀知识库下载器 GUI

<div align="center">

![语雀知识库下载器](https://img.shields.io/badge/语雀-知识库下载器-blue)
![Wails](https://img.shields.io/badge/Wails-v2-green)
![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go)
![License](https://img.shields.io/badge/license-MIT-orange)

一个现代化的语雀知识库批量下载工具,支持将语雀文档导出为 Markdown 格式,并自动下载图片到本地。

**作者: [Spritualkb](https://github.com/Spritualkb)**

</div>

## ✨ 功能特性

- 📚 **批量下载** - 支持同时管理多个下载任务
- 📊 **独立进度** - 每个任务显示独立的实时进度条
- 🎯 **任务管理** - 添加、删除、开始、暂停任务
- 📝 **批量导入** - 从文本快速导入多个下载任务
- 🖼️ **图片本地化** - 自动下载所有图片并更新为相对路径
- 🔐 **私有知识库** - 支持使用 Cookie 访问私有知识库
- ⚙️ **灵活配置** - 可自定义下载延迟、超时等参数
- 🎨 **管理后台风格** - 左右布局,清晰的任务统计和管理
- 🚀 **跨平台** - 支持 Windows、macOS、Linux

## 📸 界面预览

### 管理后台风格布局

- **左侧边栏**:
  - Logo 和作者标识
  - 实时任务统计(总任务、运行中、等待中、已完成)
  - 批量操作(开始全部、清除完成、批量导入)
  - 高级设置

- **右侧主区域**:
  - 添加任务表单
  - 任务列表卡片
  - 每个任务独立的进度条和状态

## 🚀 快速开始

### 下载使用

1. 从 [Releases](https://github.com/your-username/yuque-spider-gui/releases) 页面下载适合你系统的版本
2. 解压并运行程序
3. 在右侧表单输入语雀知识库 URL
4. 选择保存路径
5. 点击"➕ 添加任务"
6. 点击任务卡片上的 ▶️ 按钮开始下载

### 批量导入

1. 点击左侧边栏的"📝 批量导入"按钮
2. 在弹出的对话框中输入任务列表
3. 格式: 每行一个任务 `URL,Cookie(可选)`
   ```
   https://www.yuque.com/user/book1
   https://www.yuque.com/user/book2,yuque_session=xxx
   ```
4. 点击"导入"批量添加任务

### 本地开发

#### 环境要求

- Go 1.23+
- Node.js 20+
- Wails CLI v2

#### 安装依赖

```bash
# 安装 Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# 安装项目依赖
go mod tidy

# 安装前端依赖
cd frontend
npm install
cd ..
```

#### 运行开发模式

```bash
wails dev
```

#### 构建生产版本

```bash
# 构建当前平台
wails build

# 跨平台构建
wails build -platform windows/amd64
wails build -platform darwin/amd64
wails build -platform darwin/arm64
wails build -platform linux/amd64
```

## 📖 使用说明

### 基础功能

1. **添加单个任务**
   - 输入语雀知识库 URL
   - (可选) 输入 Cookie 用于访问私有知识库
   - 选择保存目录
   - 点击"➕ 添加任务"

2. **管理任务**
   - ▶️ 开始下载
   - ⏸️ 暂停下载
   - 🗑️ 删除任务

3. **查看进度**
   - 每个任务卡片显示:
     - 知识库标题
     - 当前下载的文档
     - 进度条和百分比
     - 已完成/总文档数
     - 创建和完成时间

### 批量操作

- **▶️ 开始全部**: 启动所有等待中的任务
- **🗑️ 清除完成**: 清除已完成和失败的任务
- **📝 批量导入**: 从文本批量导入多个任务

### 任务统计

左侧边栏实时显示:
- 总任务数
- 运行中任务数
- 等待中任务数
- 已完成任务数

### 高级设置

- **最小/最大延迟**: 控制下载文档之间的等待时间
- **请求超时**: 网络请求的最长等待时间

## 🛠️ 技术栈

- **后端**: Go 1.23
- **前端**: Svelte
- **框架**: Wails v2
- **依赖**:
  - goquery - HTML 解析
  - Wails Runtime - 桌面应用框架

## 📁 项目结构

```
yuque-spider-gui/
├── internal/
│   └── spider/          # 爬虫核心逻辑
│       ├── types.go     # 数据结构定义
│       ├── fetcher.go   # 网络请求处理
│       ├── downloader.go # 文档和图片下载
│       └── spider.go    # 主爬虫逻辑
├── frontend/
│   └── src/
│       └── App.svelte   # 管理后台界面
├── app.go               # 应用后端接口(任务管理)
├── main.go              # 应用入口
└── .github/
    └── workflows/
        └── build.yml    # 自动构建配置
```

## 🎯 核心特性说明

### 多任务管理

- 使用 `sync.RWMutex` 保证并发安全
- 每个任务独立的 context 和 cancel 机制
- 实时任务状态更新

### 独立进度显示

- 每个任务维护自己的进度信息
- 通过 WebSocket 事件实时推送更新
- 百分比、文档数、当前文档名等详细信息

### 批量操作

- 支持从文本快速导入多个任务
- 一键启动所有待处理任务
- 一键清理已完成任务

## 🤝 贡献

欢迎提交 Issue 和 Pull Request!

## 📝 更新日志

### v2.0.0 - 批量管理版本

#### 新增
- ✅ 批量任务管理系统
- ✅ 每个任务独立进度条
- ✅ 管理后台风格布局(左右布局)
- ✅ 实时任务统计面板
- ✅ 批量导入功能
- ✅ 批量操作(开始全部、清除完成)
- ✅ 任务状态管理(pending、running、completed、failed、cancelled)

#### 改进
- 🔄 从单任务重构为多任务架构
- 🔄 优化 UI/UX,采用管理后台设计
- 🔄 任务列表实时更新
- 🔄 更详细的任务信息展示

### v1.0.0 - 初始版本

- ✅ 使用 Go + Wails 重构
- ✅ 支持动态配置参数
- ✅ 图形化界面
- ✅ GitHub Actions 自动构建

### 基于原项目
- 原项目: [yuque-crawl](https://github.com/burpheart/yuque-crawl)

## 📄 许可证

MIT License

## 👨‍💻 作者

**Spritualkb**

- Email: 3058886310@qq.com

---

<div align="center">

**如果这个项目对你有帮助,请给一个 ⭐️ Star 支持一下!**

Made with ❤️ by Spritualkb

</div>
