<div align="center">
  <h1>简历匠</h1>
  <p>AI 驱动的中英文简历制作工具，涵盖模板、导出、分享与管理后台。</p>
  <p>
    <a href="README.md">English</a> ·
    <a href="README.zh-CN.md">简体中文</a>
  </p>
</div>

### 项目简介
- 开源全栈简历制作工具，专注高效、ATS 友好与良好体验。
- AI 辅助写作与润色。
- 通过远程 Chromium 导出 PDF/PNG。
- 支持公开分享链接与可见性管理，并提供访问数据统计。
- 提供管理后台：用户、简历、模板与分享链接。
- 内置中英文国际化。

### 主要特性
- 前端：React + Vite + TypeScript
- 后端：Go + Gin、Gorm、Redis、JWT
- 登录：GitHub / 微信 OAuth，邮箱验证码（SMTP）
- 权限：首个注册用户会被赋予 `admin` 角色；之后为 `user`。
- 生产部署：Docker 优先，`docker-compose` 一键启动
- 导出服务具备断路器以提升稳定性
- 日志与监控：Zap JSON 日志、Prometheus 指标、Loki/Grafana（可选）

### 架构说明
- 前端：React + Vite + TypeScript
- 后端：Go + Gin（包含鉴权、管理后台等能力）
- 存储：MySQL/SQLite + Redis（可选接入对象存储）
- 导出：远程 Chromium（推荐 Browserless）
- 静态资源：Nginx

### 界面与入口
- 模板、编辑器、管理、打印等页面已实现，可直接在应用中体验。

### 快速开始
- 环境要求：
  - Go 1.24+
  - Node.js 20+
  - MySQL 8.x（或 SQLite）
  - Redis 7.x
  - 可用的 Chromium DevTools（推荐 Browserless）

- 本地启动（后端）：
  - 配置环境变量（见 配置）。
  - 启动：`go run .`
  - 默认端口：`:8080`

- 初始化默认种子数据（岗位分类/岗位/内容预设）：
  - 一次性导入：`go run . seed import-default`
  - 如需导入自定义数据，可使用 `seed` 子命令从文件导入。

- 本地启动（前端）：
  - `cd frontend`
  - `npm install`
  - `npm run dev`
  - Dev 服务器默认地址 `http://localhost:3000`
  - 通过 `.env.local` 配置后端基础地址与 OAuth 弹窗通信来源。

- Docker Compose：
  - 根据需要调整 `docker-compose.yml` 的环境变量。
  - 启动：`docker compose up -d`
  - 服务包含：backend、frontend、mysql、redis、chrome。
  - 前端 Nginx 已配置后端转发与公共静态资源访问。

### 配置
后端进程环境变量（启动时读取）：
- `PORT`：HTTP 端口
- `DB_DSN`：MySQL DSN；设置后使用 MySQL
- `SQLITE_PATH`：SQLite 文件路径；当 `DB_DSN` 为空时使用
- `REDIS_ADDR`、`REDIS_PASSWORD`：Redis 连接
- `JWT_SECRET`：JWT 签名密钥

系统配置键（存储于数据库；首启可从环境变量注入；可在管理后台修改）：
- `CORS_ORIGINS`：允许的跨域来源（逗号分隔）
- `FRONTEND_BASE_URL`：前端公共地址（用于导出与 OAuth 回跳）
- `CHROME_API_URL`：Browserless/Chromium 服务的基础地址（用于 PDF/PNG 导出）
- `enable_email_verification`：注册是否启用邮箱验证码（默认 `false`）
- SMTP：`SMTP_HOST`、`SMTP_PORT`、`SMTP_USERNAME`、`SMTP_PASSWORD`、`SMTP_FROM_NAME`
- 微信 OAuth：`WECHAT_APP_ID`、`WECHAT_APP_SECRET`、`WECHAT_REDIRECT_URI`，开关 `enabled_wechat_login`
- 公众号扫码登录（关注即登录）：`WECHAT_MP_APP_ID`、`WECHAT_MP_APP_SECRET`、`WECHAT_MP_TOKEN`（可选 `WECHAT_MP_AES_KEY`），开关 `enabled_wechat_mp_login`；回调路由需公网 HTTPS 可达
- GitHub OAuth：`GITHUB_CLIENT_ID`、`GITHUB_CLIENT_SECRET`、`GITHUB_REDIRECT_URI`，开关 `enabled_github_login`
- S3 存储：`S3_BUCKET`、`S3_REGION`、`S3_ENDPOINT`、`S3_ACCESS_KEY`、`S3_SECRET_KEY`，开关 `enabled_storage_s3`

前端环境变量：
- `VITE_API_BASE`：后端基础路径或 URL
- `VITE_OAUTH_ALLOWED_ORIGINS`：允许进行 OAuth 弹窗消息通信的来源（逗号分隔）

配置示例：
- 本地后端（bash）：

```bash
export PORT=8080
export DB_DSN=""                            # 为空时使用 SQLite
export SQLITE_PATH="cvforge.db"          # SQLite 文件路径
export REDIS_ADDR="127.0.0.1:6379"
export REDIS_PASSWORD=""
export JWT_SECRET="请替换为强随机密钥"
# 首次运行用于注入系统参数（可在后台后续修改）
export CORS_ORIGINS="http://localhost:3000"
export FRONTEND_BASE_URL="http://localhost:3000"
export CHROME_API_URL="http://localhost:3000"  # Browserless/Chromium 基址，用于 PDF/PNG 导出
export SMTP_HOST=""
export SMTP_PORT=""
export SMTP_USERNAME=""
export SMTP_PASSWORD=""
export SMTP_FROM_NAME=""
export WECHAT_APP_ID=""
export WECHAT_APP_SECRET=""
export WECHAT_REDIRECT_URI=""
export WECHAT_MP_APP_ID=""
export WECHAT_MP_APP_SECRET=""
export WECHAT_MP_TOKEN=""
export WECHAT_MP_AES_KEY=""
export GITHUB_CLIENT_ID=""
export GITHUB_CLIENT_SECRET=""
export GITHUB_REDIRECT_URI=""
export S3_BUCKET=""
export S3_REGION=""
export S3_ENDPOINT=""
export S3_ACCESS_KEY=""
export S3_SECRET_KEY=""
```

- 前端 `.env.local`：

```bash
VITE_API_BASE=http://localhost:8080
VITE_OAUTH_ALLOWED_ORIGINS=http://localhost:3000
```

- Docker Compose：
  - 参考 `docker-compose.yml` 中的示例值（已为后端与前端构建参数接好）。
  - 将 `CORS_ORIGINS`、`FRONTEND_BASE_URL`、`CHROME_API_URL` 与前端的 `VITE_OAUTH_ALLOWED_ORIGINS` 替换为你的域名。

### 日志与监控
- 请求日志：Zap JSON 格式记录 HTTP 访问日志（包含请求耗时、request_id 等）。
- 请求链路：支持请求 ID 透传，便于端到端追踪。
- 指标：内置 Prometheus 指标（请求计数与耗时直方图）。
- Loki/Grafana（可选）：提供示例栈用于聚合容器日志与可视化。

### 导出服务
- 通过 DevTools WebSocket 控制远程 Chromium。
- `CHROME_API_URL` 需指向用于 PDF/PNG 导出的 Browserless/Chromium 服务。
- 服务会在导出时打开打印页面并执行 PDF 或截图生成。
- 确保 Chromium 已安装中文字体以正确显示 CJK 内容。
- 提供异步导出队列：任务由内置工作线程处理并上传存储。

### 国际化
- 内置中英文切换。
- 默认语言为英文，可在页脚或设置中切换。

### 部署建议
- 推荐使用 `docker-compose`，包含 Nginx 与 Browserless Chromium。
- 通过安全方式管理敏感变量，避免硬编码。
- 生产环境设置 `FRONTEND_BASE_URL` 为你的域名，配置好 `CORS_ORIGINS`，并将前端的 `VITE_OAUTH_ALLOWED_ORIGINS` 传入构建。

### 安全
- 妥善保管 `JWT_SECRET`、SMTP 与 OAuth 密钥。
- 严格限定 `VITE_OAUTH_ALLOWED_ORIGINS`。
- 不要将 DevTools 对外公开，建议仅内部访问。

### 参与贡献
- Fork 仓库，创建特性分支，提交 PR 并清晰描述改动。
- 保持改动内聚，必要时更新文档。
- 欢迎提 Issue 与建议。

### 许可证
- MIT 许可。参见 `LICENSE`。
