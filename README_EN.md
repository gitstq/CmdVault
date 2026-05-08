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
  <strong>Intelligent Command History Manager - Never Forget Your Terminal Commands Again</strong>
</p>

<p align="center">
  <img src="https://via.placeholder.com/800x400/7C3AED/FFFFFF?text=CmdVault+TUI+Demo" alt="CmdVault Demo" width="80%">
</p>

---

## 🎉 Introduction

**CmdVault** is an **intelligent command history manager** designed specifically for developers. It solves a pain point every developer faces: **forgetting complex commands they've used before**.

### 🎯 Problems We Solve

- 😫 **Command Amnesia**: Can't remember complex Docker, Kubernetes, or Git commands?
- 🔍 **Inefficient Search**: Traditional `Ctrl+R` only does exact matching, can't find what you need?
- 📂 **Scattered History**: Can't unify history across multiple terminal windows?
- 🔐 **Privacy Concerns**: Don't want to upload command history to the cloud?

### ✨ Unique Features

- 🧠 **Smart Fuzzy Search**: Intelligent ranking based on frequency and recency
- ⭐ **Favorites & Tags**: Bookmark frequently used commands with tags
- 📊 **Usage Statistics**: Understand your command patterns
- 🔒 **Local-First**: All data stored in local SQLite database
- 🎨 **Beautiful TUI**: Modern terminal interface built with Bubble Tea
- 🚀 **Zero Dependencies**: Single binary, works out of the box

---

## ✨ Core Features

| Feature | Description |
|---------|-------------|
| 🔍 **Fuzzy Search** | Smart matching algorithm, find commands with keywords |
| ⭐ **Favorites** | Mark frequently used commands for quick access |
| 🏷️ **Tag System** | Add tags to organize commands by project/type |
| 📝 **Notes** | Add descriptions to complex commands |
| 📊 **Statistics** | View command frequency and usage stats |
| 🔄 **History Import** | Import existing shell history with one command |
| 🖥️ **Interactive TUI** | Beautiful terminal interface with keyboard navigation |
| 🌐 **Cross-Platform** | Supports macOS, Linux, and Windows |

---

## 🚀 Quick Start

### 📋 Requirements

- **Go 1.18+** (only needed when building from source)
- **GCC** (for SQLite CGO compilation)

### 📦 Installation

#### Option 1: Download Pre-built Binary (Recommended)

```bash
# macOS / Linux
curl -sL https://github.com/gitstq/CmdVault/releases/latest/download/cmdvault-$(uname -s)-$(uname -m) -o cmdvault
chmod +x cmdvault
sudo mv cmdvault /usr/local/bin/

# Windows (PowerShell)
Invoke-WebRequest -Uri "https://github.com/gitstq/CmdVault/releases/latest/download/cmdvault-windows-amd64.exe" -OutFile "cmdvault.exe"
```

#### Option 2: Build from Source

```bash
# Clone repository
git clone https://github.com/gitstq/CmdVault.git
cd CmdVault

# Build and install
make install

# Or run directly
make build && ./dist/cmdvault
```

### 🎮 Basic Usage

```bash
# Launch interactive TUI
cmdvault search

# Search commands
cmdvault search docker

# List recent commands
cmdvault list

# Add command with tags and notes
cmdvault add "docker compose up -d" -t "docker,dev" -n "Start dev environment"

# View statistics
cmdvault stats

# Import existing history
cmdvault import

# Generate shell integration script
cmdvault init bash >> ~/.bashrc
```

---

## 📖 Detailed Usage Guide

### 🖥️ TUI Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `Tab` | Switch between search/browse mode |
| `Enter` | Select command and output |
| `f` | Toggle favorites filter |
| `s` | Star/unstar current command |
| `d` | Delete current command |
| `Esc` / `Ctrl+C` | Exit |

### 🏷️ Tag Management

```bash
# Add command with tags
cmdvault add "kubectl get pods -A" -t "k8s,monitor"

# Search by tag
cmdvault search k8s
```

### 📊 Statistics

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

### 🔄 Shell Integration

CmdVault can integrate deeply with your shell to automatically track executed commands:

```bash
# Bash
echo 'eval "$(cmdvault init bash)"' >> ~/.bashrc

# Zsh
echo 'eval "$(cmdvault init zsh)"' >> ~/.zshrc

# Fish
cmdvault init fish > ~/.config/fish/conf.d/cmdvault.fish
```

---

## 💡 Design Philosophy & Roadmap

### 🎨 Design Principles

1. **Local-First**: Data stored entirely locally, protecting user privacy
2. **Simple & Efficient**: Single responsibility, focused on command management
3. **Developer-Friendly**: Keyboard-driven, matching developer habits
4. **Extensible**: Modular design for easy feature additions

### 🗓️ Roadmap

- [ ] **v1.1** - AI-powered semantic search (optional)
- [ ] **v1.2** - Cross-device sync (end-to-end encrypted)
- [ ] **v1.3** - Command templates with variable substitution
- [ ] **v1.4** - Team shared command library
- [ ] **v1.5** - Web UI interface

---

## 📦 Build & Deployment

### Local Build

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test

# Lint code
make lint
```

### Cross-Platform Binaries

| Platform | Architecture | Filename |
|----------|--------------|----------|
| macOS | amd64 | `cmdvault-darwin-amd64` |
| macOS | arm64 | `cmdvault-darwin-arm64` |
| Linux | amd64 | `cmdvault-linux-amd64` |
| Linux | arm64 | `cmdvault-linux-arm64` |
| Windows | amd64 | `cmdvault-windows-amd64.exe` |
| Windows | arm64 | `cmdvault-windows-arm64.exe` |

---

## 🤝 Contributing

We welcome all forms of contributions!

### How to Contribute

1. Fork this repository
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'feat: Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Create a Pull Request

### Commit Convention

Please follow [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation update
- `refactor:` Code refactoring
- `test:` Test related

---

## 📄 License

This project is licensed under the **MIT License** - you are free to use, modify, and distribute this software.

See [LICENSE](LICENSE) file for details.

---

<p align="center">
  Made with ❤️ by <a href="https://github.com/gitstq">gitstq</a>
</p>

<p align="center">
  If this project helps you, please give it a ⭐ Star!
</p>
