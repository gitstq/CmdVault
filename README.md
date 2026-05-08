<p align="center">
  <img src="https://img.shields.io/badge/version-1.0.0-blue.svg" alt="Version">
  <img src="https://img.shields.io/badge/license-MIT-green.svg" alt="License">
  <img src="https://img.shields.io/badge/platform-macOS%20%7C%20Linux%20%7C%20Windows-orange.svg" alt="Platform">
  <img src="https://img.shields.io/badge/Go-1.18+-00ADD8.svg" alt="Go Version">
</p>

<p align="center">
  <a href="README.md">简体中文</a> | <a href="README_EN.md">English</a> | <a href="README_TW.md">繁體中文</a>
</p>

<h1 align="center">🔐 CmdVault</h1>

<p align="center">
  <strong>智能命令历史管理器 - 让你的终端命令永不遗忘</strong>
</p>

<p align="center">
  <img src="https://via.placeholder.com/800x400/7C3AED/FFFFFF?text=CmdVault+TUI+Demo" alt="CmdVault Demo" width="80%">
</p>

---

## 🎉 项目介绍

**CmdVault** 是一款专为开发者打造的**智能命令历史管理工具**。它解决了每个开发者都会遇到的痛点：**忘记曾经使用过的复杂命令**。

### 🎯 解决的痛点

- 😫 **命令遗忘**：复杂的 Docker、Kubernetes、Git 命令总是记不住？
- 🔍 **搜索低效**：传统的 `Ctrl+R` 只能精确匹配，找不到想要的命令？
- 📂 **历史分散**：多个终端窗口的历史记录无法统一管理？
- 🔐 **隐私担忧**：不想把命令历史上传到云端？

### ✨ 自研差异化亮点

- 🧠 **智能模糊搜索**：基于使用频率和时间的智能排序，输入关键词即可快速定位
- ⭐ **收藏与标签**：将常用命令收藏、打标签，一键快速调用
- 📊 **使用统计**：了解你的命令使用习惯，发现高频操作
- 🔒 **本地优先**：所有数据存储在本地 SQLite 数据库，隐私无忧
- 🎨 **精美 TUI**：基于 Bubble Tea 构建的现代化终端界面
- 🚀 **零依赖**：单二进制文件，开箱即用

---

## ✨ 核心特性

| 特性 | 描述 |
|------|------|
| 🔍 **模糊搜索** | 智能匹配算法，输入关键词即可快速定位命令 |
| ⭐ **收藏功能** | 将常用命令标记为收藏，快速访问 |
| 🏷️ **标签系统** | 为命令添加标签，按项目/类型分类管理 |
| 📝 **备注功能** | 为复杂命令添加使用说明 |
| 📊 **统计分析** | 查看命令使用频率、执行次数等统计 |
| 🔄 **历史导入** | 一键导入现有 Shell 历史记录 |
| 🖥️ **交互式 TUI** | 美观的终端用户界面，键盘操作流畅 |
| 🌐 **跨平台** | 支持 macOS、Linux、Windows |

---

## 🚀 快速开始

### 📋 环境要求

- **Go 1.18+**（仅从源码编译时需要）
- **GCC**（用于编译 SQLite CGO）

### 📦 安装方式

#### 方式一：下载预编译二进制（推荐）

```bash
# macOS / Linux
curl -sL https://github.com/gitstq/CmdVault/releases/latest/download/cmdvault-$(uname -s)-$(uname -m) -o cmdvault
chmod +x cmdvault
sudo mv cmdvault /usr/local/bin/

# Windows (PowerShell)
Invoke-WebRequest -Uri "https://github.com/gitstq/CmdVault/releases/latest/download/cmdvault-windows-amd64.exe" -OutFile "cmdvault.exe"
```

#### 方式二：从源码编译

```bash
# 克隆仓库
git clone https://github.com/gitstq/CmdVault.git
cd CmdVault

# 编译安装
make install

# 或者直接运行
make build && ./dist/cmdvault
```

### 🎮 基本使用

```bash
# 启动交互式 TUI 界面
cmdvault search

# 搜索命令
cmdvault search docker

# 列出最近命令
cmdvault list

# 添加命令（带标签和备注）
cmdvault add "docker compose up -d" -t "docker,dev" -n "启动开发环境"

# 查看统计
cmdvault stats

# 导入现有历史
cmdvault import

# 生成 Shell 集成脚本
cmdvault init bash >> ~/.bashrc
```

---

## 📖 详细使用指南

### 🖥️ TUI 快捷键

| 快捷键 | 功能 |
|--------|------|
| `Tab` | 切换搜索/浏览模式 |
| `Enter` | 选择命令并输出 |
| `f` | 切换收藏过滤 |
| `s` | 收藏/取消收藏当前命令 |
| `d` | 删除当前命令 |
| `Esc` / `Ctrl+C` | 退出 |

### 🏷️ 标签管理

```bash
# 添加带标签的命令
cmdvault add "kubectl get pods -A" -t "k8s,monitor"

# 搜索特定标签
cmdvault search k8s
```

### 📊 统计信息

```bash
$ cmdvault stats

📊 CmdVault Statistics
━━━━━━━━━━━━━━━━━━━━━━
📝 Total Commands: 1,234
🎯 Unique Commands: 567
⭐ Favorites: 42

🏆 Top Commands:
  1. git status (156 times)
  2. docker ps (89 times)
  3. kubectl get pods (67 times)
```

### 🔄 Shell 集成

CmdVault 可以与你的 Shell 深度集成，自动记录执行的命令：

```bash
# Bash
echo 'eval "$(cmdvault init bash)"' >> ~/.bashrc

# Zsh
echo 'eval "$(cmdvault init zsh)"' >> ~/.zshrc

# Fish
cmdvault init fish > ~/.config/fish/conf.d/cmdvault.fish
```

---

## 💡 设计思路与迭代规划

### 🎨 设计理念

CmdVault 的设计遵循以下原则：

1. **本地优先**：数据完全存储在本地，保护用户隐私
2. **简单高效**：单一职责，专注于命令管理
3. **开发者友好**：键盘操作为主，符合开发者习惯
4. **可扩展**：模块化设计，便于添加新功能

### 🗓️ 迭代规划

- [ ] **v1.1** - AI 语义搜索（可选）
- [ ] **v1.2** - 跨设备同步（端到端加密）
- [ ] **v1.3** - 命令模板与变量替换
- [ ] **v1.4** - 团队共享命令库
- [ ] **v1.5** - Web UI 界面

---

## 📦 打包与部署指南

### 本地编译

```bash
# 编译当前平台
make build

# 编译所有平台
make build-all

# 运行测试
make test

# 代码检查
make lint
```

### 跨平台打包产物

| 平台 | 架构 | 文件名 |
|------|------|--------|
| macOS | amd64 | `cmdvault-darwin-amd64` |
| macOS | arm64 | `cmdvault-darwin-arm64` |
| Linux | amd64 | `cmdvault-linux-amd64` |
| Linux | arm64 | `cmdvault-linux-arm64` |
| Windows | amd64 | `cmdvault-windows-amd64.exe` |
| Windows | arm64 | `cmdvault-windows-arm64.exe` |

---

## 🤝 贡献指南

我们欢迎所有形式的贡献！

### 如何贡献

1. Fork 本仓库
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'feat: Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

### 提交规范

请遵循 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

- `feat:` 新功能
- `fix:` 修复 Bug
- `docs:` 文档更新
- `refactor:` 代码重构
- `test:` 测试相关

---

## 📄 开源协议说明

本项目采用 **MIT 协议** 开源，你可以自由地使用、修改和分发本软件。

详见 [LICENSE](LICENSE) 文件。

---

<p align="center">
  Made with ❤️ by <a href="https://github.com/gitstq">gitstq</a>
</p>

<p align="center">
  如果这个项目对你有帮助，请给一个 ⭐ Star 支持一下！
</p>
