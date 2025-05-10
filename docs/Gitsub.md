# Gitsub MVP开发文档（基础版）

## 版本：v0.1 MVP

**目标：** 用最少的时间实现核心功能，验证技术可行性

---

## 1. 项目概述

### 1.1 MVP目标

实现一个命令行工具，能够从大型GitHub仓库中只下载指定目录，而不是整个仓库。

### 1.2 核心功能（MVP范围）

* ✅ 克隆指定目录
* ✅ 支持单个或多个目录
* ✅ 基本的进度显示
* ❌ 交互式界面（v2.0）
* ❌ 图形
* ❌ 配置文件管理（v1.0）

---

## 2. 技术方案（极简版）

### 2.1 核心技术

```bash
# Git命令组合（这就是全部魔法）
git clone --filter=blob:none --sparse <repo-url>
cd <repo>
git sparse-checkout set <directory>
```

### 2.2 工作原理

1. `--filter=blob:none`：只下载commit和tree对象，不下载文件内容
2. `--sparse`：启用稀疏检出
3. `sparse-checkout set`：指定要检出的目录
4. Git自动下载指定目录的文件内容

---

## 3. MVP功能规格

### 3.1 命令格式

```bash
# 基础用法
gitsub clone <repo-url> <directory>

# 示例
gitsub clone https://github.com/tensorflow/tensorflow tensorflow/core

# 多目录
gitsub clone https://github.com/user/repo dir1 dir2 dir3
```

### 3.2 参数说明

````
gitsub clone <repository> <directory> [directory...]

参数：
  repository    GitHub仓库URL
  directory     要克隆的目录路径（一个或多个）

选项：
  -b, --branch <branch>    指定分支（默认：main）
  -o, --output <path>      输出目录（默认：当前目录/术栈（推荐）

### 4.1 语言选择：Go
**理由：**
- 单文件分发（无依赖）
- 跨平台编译简单
- 性能好，启动快

### 4.2 依赖库（最小化）
```go
// 只需标准库 + 1个CLI库
import (
    "os//spf13/cobra"  // CLI框架（可选，也可纯标准库）
)
````

---

## 5. 项目结构（极简）

```
gitsub/
├── main.go              # 入口文件（200行左右）
├── clone.go             # 克隆逻辑（100行）
├── git.go               # Git命令封装（50行）
├── utils.go             # 工具函数（50行）
├── go.mod
└── README.md

总代码量：~400行
```

---

## 6. 核心代码设计

### 6.1 main.go（入口）

```go
package main

import (
    "flag"
    "fmt"
    "os"
)

func main() {
    // 解析参数
    branch := flag.String("b", "", "分支名称")
    output := flag.String
    
    args := flag.Args()
    if len(args) < 2 {
        fmt.Println("用法: gitsub clone <repo-url> <directory> [directory...]")
        os.Exit(1)
    }
    
    repoURL := args[0]
    directories := args[1:]
    
    // 执行克隆
    err := Clone(repoURL, directories, *branch, *output)
    if err != nil {
        fmt.Printf("错误: %v\n", err)
        os.Exit(1)
    }
    
    fmt.Println("✓ 克隆完成！")
}
```

### 6.2 clone.go（核心逻辑）

```go
package main

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
)

func Clone(repoURL string, directories []string, branch, output string) error {
    // 1. 提取仓库名
    repoName := extractRepoName(repoURL)
    if output == "" {
        output = repoName
    }
    
    // 2. 确定分支
    if branch == "" {
        branch = "main"  // 默认分支
    }
    
    // 3. 执行克隆
    fmt.Printf("克隆仓库: %s\n", repoURL 1: 初始化
    if err := RunGit("init", output); err != nil {
        return err
    }
    
    // Step 2: 添加remote
    if err := RunGitInDir(output, "remote", "add", "origin", repoURL); err != nil {
        return err
    }
    
    // Step 3: 配置sparse-checkout
    if err := RunGitInDir(output, "config", "core.sparseCheckout", "true"); err != nil {
        return err
    }
    
    // Step 4: 设置sparse-checkout路径
    args := append([]string{"sparse-checkout", "set", "--cone"}, directories...)
    if err := RunGitInDir(output, args...); err != nil {
        return err
    }
    
    // Step 5: Fetch
    fmt.Println("下载中...")
    if err := RunGitInDir(output, "fetch", "--filter=blob:none", "--depth=1", "origin", branch); err != nil {
        return err
    }
    
    // Step 6: Checkout
    if err := RunGitInDir(output, "checkout", branch); err != nil {
        return err
    }
    
    return nil
}

// 从URL提取仓库名
func extractRepoName(url string) string {
    // https://github.com/user/repo.git -> repo
    parts := strings.Split(strings.TrimSuffix(url, ".git"), "/")
    return parts[len(parts)-1]
}
```

### 6.3 git.go（Git命令封装）

```go
package main

import (
    "os/exec"
)

// 执行Git命令
func RunGit(args ...string) error {
    cmd := exec.Command("git", args...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    return cmd.Run()
}

// 在指定目录执行Git命令
func RunGitInDir(dir string, args ...string) error {
    cmd := exec.Command("git", args...)
    cmd.Dir = dir
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    return cmd.Run()
}

// 检查Git版本
func CheckGitVersion() error {
    cmd := exec.Command("git", "--version")
    output, err := cmd.Output()
    if err != nil {
        return fmt.Errorf("Git未安装")
    }
    
    // 简单检查（可选）
    version := string(output)
    if !strings.Contains(version, "git version") {
        return fmt.Errorf("Git版本检测失败")
    }
    
    return nil
}
```

### 6.4 utils.go（工具函数）

```go
package main

import (
    "fmt"
    "strings"
)

// 验证URL格式
func ValidateRepoURL(url string) error {
    if !strings.HasPrefix(url, "https://") && !strings.HasPrefix(url, "git@") {
        return fmt.Errorf("无效的仓库URL")
    }
    return nil
}

// 标准化目录路径
func NormalizePath(path string) string {
    // 移除前后的斜杠
    path = strings.Trim(path, "/")
    // 统一使用正斜杠
    path = strings.ReplaceAll(path, "\\", "/")
    return path
}
```

---

## 7. 开发步骤（1-2周）

### 第1天：环境准备

```bash
# 创建项目
mkdir gitsub && cd gitsub
go mod init github.com/yourusername/gitsub

# 创建文件
touch main.go clone.go git.go utils.go
```

### 第2-3天：核心功能开发

1. 实现`git.go`中的Git命令封装
2. 实现`clone.go`中的克隆逻辑
3. 实现`main.go`中的命令行解析

### 第4天：基础测试

```bash
# 手动测试
go run . clone https://github.com/tensorflow/tensorflow tensorflow/core

# 测试不同场景
go run . clone <repo> dir1 dir2
go run . clone <repo> -b develop dir1
```

### 第5天：错误处理优化

* 添加基本的错误提示
* 验证输入参数
* 检查Git版本

### 第6-7天：打包发布

```bash
# 编译
go build -o gitsub

# 跨平台编译
GOOS=linux GOARCH=amd64 go build -o gitsub-linux
GOOS=darwin GOARCH=amd64 go build -o gitsub-macos
GOOS=windows GOARCH=amd64 go build -o gitsub.exe
```

---

## 8. 测试清单（MVP）

### 8.1 基础测试

```bash
# 测试1：克隆单个目录
./gitsub clone https://github.com/kubernetes/kubernetes pkg/api

# 测试2：克隆多个目录
./gitsub clone https://github.com/user/repo src/app src/utils

# 测试3：指定分支
./gitsub clone https://github.com/user/repo -b develop src

# 测试4：指定输出目录
./gitsub clone https://github.com/user/repo -o myrepo src
```

### 8.2 错误测试

```bash
# 无效的URL
./gitsub clone invalid-url dir1

# 不存在的仓库
./gitsub clone https://github.com/user/nonexist dir1

# 不存在的目录（应该在clone后发现）
./gitsub clone https://github.com/user/repo nonexist-dir
```

---

## 9. MVP局限性（已知问题）

### 9.1 不支持的功能

* ❌ 交互式目录浏览
* ❌ GitHub API预览
* ❌ 进度百分比显示
* ❌ 配置文件管理
* ❌ Token认证
* ❌ 已克隆仓库的管理（add/remove）

### 9.2 已知限制

* 不验证目录是否存在（Git会在checkout时报错）
* 不检查网络连接
* 错误提示较简单
* 不支持断点续传

### 9.3 用户需要

* Git 2.25+已安装
* 网络连接
* 对公开仓库的访问权限

---

## 10. 使用文档（MVP）

### 10.1 安装

```bash
# 从源码编译
git clone https://github.com/yourusername/gitsub
cd gitsub
go build -o gitsub
sudo mv gitsub /usr/local/bin/
```

### 10.2 快速开始

```bash
# 克隆TensorFlow的core目录
gitsub clone https://github.com/tensorflow/tensorflow tensorflow/core

# 进入目录查看
cd tensorflow
ls tensorflow/core
```

### 10.3 常见用法

```bash
# 克隆前端代码
gitsub clone https://github.com/facebook/react packages/react packages/react-dom

# 克隆文档
gitsub clone https://github.com/user/repo docs

# 克隆到指定目录
gitsub clone https://github.com/user/repo -o myproject src/app
```

---

## 11. README.md模板

````markdown
# Gitsub

从大型GitHub仓库中只克隆你需要的目录。

## 特性

- 🚀 只下载指定目录，节省90%+的时间和空间
- 📦 单文件工具，无需依赖
- 💻 支持Linux、macOS、Windows

## 安装

```bash
# 下载预编译版本
wget https://github.com/user/gitsub/releases/latest/download/gitsub-linux
chmod +x gitsub-linux
sudo mv gitsub-linux /usr/local/bin/gitsub
```

## 使用

```bash
# 基础用法
gitsub clone <仓库URL> <目录>

# 示例：只克隆TensorFlow的core目录
gitsub clone https://github.com/tensorflow/tensorflow tensorflow/core

# 克隆多个目录
gitsub clone https://github.com/user/repo src/app src/utils docs
```

## 要求

- Git 2.25+

## 工作原理

Gitsub使用Git的Partial Clone和Sparse Checkout功能：
```bash
git clone --filter=blob:none --sparse <repo>
git sparse-checkout set <directory>
```

## 许可证

MIT
````

---

## 12. 下一步计划（v1.0）

完成MVP后，可以考虑添加：

1. **基础管理功能**

   * `gitsub add <dir>` - 添加目录
   * `gitsub list` - 列出当前配置格式验证

2. **简单进度显示**

   * 显示当前执行的步骤

3. **基础文档**

   * 完善README
   * 添加示例

---

## 总结

### MVP开发时间估算

* **核心功能开发**: 3-4天
* **测试调试**: 2-3天
* **文档编写**: 1天
* **总计**: **1周**

### MVP代码量

* 总代码：~400行Go代码
* 复杂度：低（主要是调用Git命令）

### MVP验证目标

* ✅ 技术可行性验证
* ✅ 用户需求验证
* ✅ 获得早期反馈
* ✅ 为后续开发奠定基础

---

**开始开发吧！先把MVP做出来，再考虑高级功能！** 🚀
