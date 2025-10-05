# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added
- ✨ 使用 Go + Wails 重构的现代化 GUI 应用
- 🎨 美观的渐变色界面设计
- 📊 实时下载进度显示
- ⚙️ 可配置的下载参数(延迟、超时等)
- 🔄 支持取消正在进行的下载
- 📁 可视化目录选择器
- 🌐 跨平台支持(Windows/macOS/Linux)
- 🤖 GitHub Actions 自动构建和发布

### Changed
- 🔄 从 Python CLI 迁移到 Go GUI
- ♻️ 重构代码架构,模块化设计
- 📝 去除硬编码配置,全部可自定义

### Fixed
- 🐛 优化图片下载逻辑
- 🔧 改进错误处理和用户提示

## [原项目] - Python CLI 版本

基于 [yuque-crawl](https://github.com/burpheart/yuque-crawl) 项目

### Features
- 下载语雀知识库为 Markdown
- 支持 Cookie 认证
- 图片本地化
- 生成目录索引
