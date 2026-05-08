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
  <strong>智慧型指令歷史管理器 - 讓你的終端指令永不遺忘</strong>
</p>

<p align="center">
  <img src="https://via.placeholder.com/800x400/7C3AED/FFFFFF?text=CmdVault+TUI+Demo" alt="CmdVault Demo" width="80%">
</p>

---

## 🎉 專案介紹

**CmdVault** 是一款專為開發者打造的**智慧型指令歷史管理工具**。它解決了每個開發者都會遇到的痛點：**忘記曾經使用過的複雜指令**。

### 🎯 解決的痛點

- 😫 **指令遺忘**：複雜的 Docker、Kubernetes、Git 指令總是記不住？
- 🔍 **搜尋低效**：傳統的 `Ctrl+R` 只能精確匹配，找不到想要的指令？
- 📂 **歷史分散**：多個終端視窗的歷史記錄無法統一管理？
- 🔐 **隱私擔憂**：不想把指令歷史上傳到雲端？

### ✨ 自研差異化亮點

- 🧠 **智慧模糊搜尋**：基於使用頻率和時間的智慧排序，輸入關鍵字即可快速定位
- ⭐ **收藏與標籤**：將常用指令收藏、打標籤，一鍵快速調用
- 📊 **使用統計**：了解你的指令使用習慣，發現高頻操作
- 🔒 **本地優先**：所有資料儲存在本地 SQLite 資料庫，隱私無憂
- 🎨 **精美 TUI**：基於 Bubble Tea 構建的現代化終端介面
- 🚀 **零依賴**：單一二進位檔案，開箱即用

---

## ✨ 核心特性

| 特性 | 描述 |
|------|------|
| 🔍 **模糊搜尋** | 智慧匹配演算法，輸入關鍵字即可快速定位指令 |
| ⭐ **收藏功能** | 將常用指令標記為收藏，快速存取 |
| 🏷️ **標籤系統** | 為指令添加標籤，按專案/類型分類管理 |
| 📝 **備註功能** | 為複雜指令添加使用說明 |
| 📊 **統計分析** | 查看指令使用頻率、執行次數等統計 |
| 🔄 **歷史匯入** | 一鍵匯入現有 Shell 歷史記錄 |
| 🖥️ **互動式 TUI** | 美觀的終端使用者介面，鍵盤操作流暢 |
| 🌐 **跨平台** | 支援 macOS、Linux、Windows |

---

## 🚀 快速開始

### 📋 環境要求

- **Go 1.18+**（僅從原始碼編譯時需要）
- **GCC**（用於編譯 SQLite CGO）

### 📦 安裝方式

#### 方式一：下載預編譯二進位（推薦）

```bash
# macOS / Linux
curl -sL https://github.com/gitstq/CmdVault/releases/latest/download/cmdvault-$(uname -s)-$(uname -m) -o cmdvault
chmod +x cmdvault
sudo mv cmdvault /usr/local/bin/

# Windows (PowerShell)
Invoke-WebRequest -Uri "https://github.com/gitstq/CmdVault/releases/latest/download/cmdvault-windows-amd64.exe" -OutFile "cmdvault.exe"
```

#### 方式二：從原始碼編譯

```bash
# 複製儲存庫
git clone https://github.com/gitstq/CmdVault.git
cd CmdVault

# 編譯安裝
make install

# 或者直接執行
make build && ./dist/cmdvault
```

### 🎮 基本使用

```bash
# 啟動互動式 TUI 介面
cmdvault search

# 搜尋指令
cmdvault search docker

# 列出最近指令
cmdvault list

# 新增指令（帶標籤和備註）
cmdvault add "docker compose up -d" -t "docker,dev" -n "啟動開發環境"

# 查看統計
cmdvault stats

# 匯入現有歷史
cmdvault import

# 產生 Shell 整合腳本
cmdvault init bash >> ~/.bashrc
```

---

## 📖 詳細使用指南

### 🖥️ TUI 快速鍵

| 快速鍵 | 功能 |
|--------|------|
| `Tab` | 切換搜尋/瀏覽模式 |
| `Enter` | 選擇指令並輸出 |
| `f` | 切換收藏過濾 |
| `s` | 收藏/取消收藏當前指令 |
| `d` | 刪除當前指令 |
| `Esc` / `Ctrl+C` | 離開 |

### 🏷️ 標籤管理

```bash
# 新增帶標籤的指令
cmdvault add "kubectl get pods -A" -t "k8s,monitor"

# 搜尋特定標籤
cmdvault search k8s
```

### 📊 統計資訊

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

### 🔄 Shell 整合

CmdVault 可以與你的 Shell 深度整合，自動記錄執行的指令：

```bash
# Bash
echo 'eval "$(cmdvault init bash)"' >> ~/.bashrc

# Zsh
echo 'eval "$(cmdvault init zsh)"' >> ~/.zshrc

# Fish
cmdvault init fish > ~/.config/fish/conf.d/cmdvault.fish
```

---

## 💡 設計思路與迭代規劃

### 🎨 設計理念

CmdVault 的設計遵循以下原則：

1. **本地優先**：資料完全儲存在本地，保護使用者隱私
2. **簡單高效**：單一職責，專注於指令管理
3. **開發者友善**：鍵盤操作為主，符合開發者習慣
4. **可擴展**：模組化設計，便於新增功能

### 🗓️ 迭代規劃

- [ ] **v1.1** - AI 語意搜尋（可選）
- [ ] **v1.2** - 跨裝置同步（端對端加密）
- [ ] **v1.3** - 指令範本與變數替換
- [ ] **v1.4** - 團隊共享指令庫
- [ ] **v1.5** - Web UI 介面

---

## 📦 打包與部署指南

### 本地編譯

```bash
# 編譯當前平台
make build

# 編譯所有平台
make build-all

# 執行測試
make test

# 程式碼檢查
make lint
```

### 跨平台打包產物

| 平台 | 架構 | 檔案名稱 |
|------|------|----------|
| macOS | amd64 | `cmdvault-darwin-amd64` |
| macOS | arm64 | `cmdvault-darwin-arm64` |
| Linux | amd64 | `cmdvault-linux-amd64` |
| Linux | arm64 | `cmdvault-linux-arm64` |
| Windows | amd64 | `cmdvault-windows-amd64.exe` |
| Windows | arm64 | `cmdvault-windows-arm64.exe` |

---

## 🤝 貢獻指南

我們歡迎所有形式的貢獻！

### 如何貢獻

1. Fork 本儲存庫
2. 建立功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交變更 (`git commit -m 'feat: Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 建立 Pull Request

### 提交規範

請遵循 [Conventional Commits](https://www.conventionalcommits.org/) 規範：

- `feat:` 新功能
- `fix:` 修復 Bug
- `docs:` 文件更新
- `refactor:` 程式碼重構
- `test:` 測試相關

---

## 📄 開源協議說明

本專案採用 **MIT 協議** 開源，你可以自由地使用、修改和分發本軟體。

詳見 [LICENSE](LICENSE) 文件。

---

<p align="center">
  Made with ❤️ by <a href="https://github.com/gitstq">gitstq</a>
</p>

<p align="center">
  如果這個專案對你有幫助，請給一個 ⭐ Star 支持一下！
</p>
