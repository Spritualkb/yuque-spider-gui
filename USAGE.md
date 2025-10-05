# 使用指南

## 安装与运行

### Windows 用户

1. 下载 `yuque-spider-gui-vX.X.X-windows-amd64.zip`
2. 解压到任意目录
3. 双击 `yuque-spider-gui.exe` 运行

### macOS 用户

1. 下载对应芯片的版本:
   - Intel 芯片: `yuque-spider-gui-vX.X.X-darwin-amd64.tar.gz`
   - Apple Silicon: `yuque-spider-gui-vX.X.X-darwin-arm64.tar.gz`

2. 解压文件:
   ```bash
   tar -xzf yuque-spider-gui-vX.X.X-darwin-*.tar.gz
   ```

3. 运行应用:
   - 双击 `yuque-spider-gui.app`
   - 或命令行: `open yuque-spider-gui.app`

4. 如果遇到"无法打开,因为无法验证开发者"的提示:
   ```bash
   xattr -cr yuque-spider-gui.app
   ```

### Linux 用户

1. 下载 `yuque-spider-gui-vX.X.X-linux-amd64.tar.gz`
2. 解压并运行:
   ```bash
   tar -xzf yuque-spider-gui-vX.X.X-linux-amd64.tar.gz
   chmod +x yuque-spider-gui
   ./yuque-spider-gui
   ```

## 基础使用

### 1. 获取语雀知识库 URL

访问你想要下载的语雀知识库,复制浏览器地址栏的 URL。

例如: `https://www.yuque.com/username/bookname`

### 2. 获取 Cookie (仅私有知识库需要)

如果知识库是私有的,需要登录语雀后获取 Cookie:

1. 在 Chrome/Edge 浏览器中访问语雀
2. 按 `F12` 打开开发者工具
3. 切换到 `Application` (或 `应用`) 标签
4. 左侧找到 `Cookies` → `https://www.yuque.com`
5. 复制所有 Cookie 值,格式类似:
   ```
   _yuque_session=xxx; yuque_ctoken=yyy; ...
   ```

### 3. 开始下载

1. 在程序中粘贴 URL
2. 如有需要,粘贴 Cookie
3. 点击"选择目录"选择保存位置
4. 点击"🚀 开始下载"

### 4. 查看进度

程序会显示:
- 当前知识库名称
- 正在下载的文档名称
- 进度条和百分比
- 已完成/总文档数

## 高级设置说明

### 最小延迟和最大延迟

- **作用**: 控制下载每个文档之间的等待时间
- **建议值**: 1-4 秒
- **说明**:
  - 设置延迟可以避免请求过快被语雀限流
  - 如果遇到频繁失败,可以适当增加延迟
  - 私有知识库建议设置较长延迟(3-6秒)

### 请求超时

- **作用**: 单个网络请求的最长等待时间
- **建议值**: 30 秒
- **说明**:
  - 网络较慢时可以适当增加
  - 一般不需要修改默认值

## 常见问题

### Q: 下载失败怎么办?

**A:** 可能的原因和解决方法:

1. **URL 错误** - 检查 URL 是否完整正确
2. **Cookie 过期** - 重新获取 Cookie
3. **网络问题** - 检查网络连接
4. **请求限流** - 增加延迟时间

### Q: 部分图片无法显示?

**A:**
- 程序会自动下载所有图片到 `assets` 目录
- 如果某些图片下载失败,会保留原始链接
- 可以手动访问原始链接下载

### Q: 如何导入到其他笔记软件?

**A:** 下载的文件是标准 Markdown 格式:

- **Obsidian**: 直接将文件夹作为 vault 打开
- **Notion**: 使用 Notion 的导入功能
- **Typora**: 直接打开 Markdown 文件
- **语雀**: 可以重新导入

### Q: SUMMARY.md 是什么?

**A:**
- 自动生成的目录索引文件
- 包含所有文档的链接
- 可以用于快速导航

### Q: 支持批量下载多个知识库吗?

**A:**
- 当前版本需要逐个下载
- 可以在一个任务完成后继续下载下一个

## 技术支持

- GitHub Issues: [提交问题](https://github.com/your-username/yuque-spider-gui/issues)
- Email: 3058886310@qq.com

## 更新日志

查看 [CHANGELOG.md](CHANGELOG.md) 了解版本更新历史。
