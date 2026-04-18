# AGENTS.md

AI 编码助手在本项目中的统一技术规范。

## 项目概述

Go + React 前后端分离但最终编译为单个二进制分发的 Web 应用模板。后端将前端构建产物通过 `//go:embed` 嵌入到 Go 二进制中，生产部署只需一个可执行文件加数据库。

## 技术栈

- 后端：Go 1.25 + `fox-gonic/fox` + GORM + PostgreSQL（默认）/ MySQL
- 前端：React 18 + TypeScript 5 + Vite 6 + Tailwind CSS 4 + shadcn/ui
  - React Router v7（路由）
  - React Query v5（服务端状态）
  - Axios（HTTP 客户端）
  - Lucide React（图标）

## 开发命令

```bash
task install        # 安装前后端依赖
task dev            # 启动开发环境（热重载）
task build          # 构建生产二进制（含前端）
task build-all      # 多平台构建
task run            # 生产模式运行
task lint           # 自动修复代码风格并执行检查
task check          # 完整检查（后端 + 前端类型 + mod tidy）
task test           # 运行测试（race 检测 + 覆盖率）
task clean          # 清理构建产物
task update-tools   # 更新开发工具
```

## 目录概览

```text
cmd/app/                      # 主程序入口与本地配置
internal/config/              # YAML 基础配置加载（支持 PostgreSQL / MySQL）
internal/entity/              # 数据模型与领域类型
internal/handler/             # HTTP 处理器、路由注册、中间件
internal/service/             # 业务逻辑、数据库操作
internal/errors/              # 统一错误类型
internal/website/             # SPA 嵌入
  ├── assets_development.go   #   开发模式：反向代理到 Vite dev server
  ├── assets_production.go    #   生产模式：go:embed 嵌入静态资源
  ├── package.json
  ├── vite.config.ts
  ├── tsconfig*.json
  ├── eslint.config.js
  ├── vitest.config.ts
  ├── components.json         #   shadcn 配置
  ├── index.html
  ├── public/
  ├── build/                  #   vite 产物（被 embed）
  └── src/
      ├── main.tsx
      ├── App.tsx
      ├── router.tsx
      ├── globals.css
      ├── api/
      ├── types/
      ├── views/
      ├── components/ (含 ui/ by shadcn)
      ├── layouts/
      ├── hooks/
      ├── context/
      └── lib/
pkg/gormlog/                  # GORM 日志适配器
```

## 核心架构约束

### 后端

- 分层保持为 `Handler -> Service -> Entity`
- 路由统一在 `internal/handler/handler.go` 中注册
- 支持 PostgreSQL（默认）和 MySQL，通过 YAML 配置 `driver` 切换
- YAML 仅保留基础启动配置（地址、数据库驱动和连接字符串）

### 前端

- 路由使用 React Router v7
- 服务端状态管理使用 React Query（`@tanstack/react-query`）
- API 调用放在 `internal/website/src/api/`
- 类型定义放在 `internal/website/src/types/`
- 页面放在 `internal/website/src/views/`
- 优先复用现有 shadcn/ui 与 Tailwind 风格

### 单二进制嵌入机制

- `internal/website/assets_development.go`（`//go:build development`）反向代理到 Vite dev server
- `internal/website/assets_production.go`（`//go:build !development`）使用 `//go:embed build/*` 嵌入静态资源
- NotFound 处理器：`/api` 前缀返回 JSON 404，其余走 SPA index fallback

## 强制规则

- 遵守现有分层和目录结构，不为局部改动重塑架构
- 提交前执行 `task check`

## 提交前最小检查

- 执行 `task check`，未通过不得提交
- 检查是否需要同步更新前端 API/types
