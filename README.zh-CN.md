<div align="center">
  <h1>OpenResume</h1>
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
- 支持公开分享链接与可见性管理。
- 提供管理后台：用户、简历、模板与分享链接。
- 内置中英文国际化。

### 主要特性
- 前端：React + Vite + TypeScript
- 后端：Go + Gin、Gorm、Redis、JWT
- 登录：GitHub / 微信 OAuth，邮箱验证码（SMTP）
- 生产部署：Docker 优先，`docker-compose` 一键启动
- 导出服务具备断路器以提升稳定性
- 日志与监控：Zap JSON 日志、Prometheus 指标、Loki/Grafana（可选）

### 架构说明
- 前端：`frontend/`（React 19, Vite 6, TypeScript）
- 后端：`internal/` 模块（Gin 处理器、服务、中间件），入口 `main.go`
- 存储：MySQL/SQLite + Redis
- 导出：通过 DevTools WebSocket 连接远程 Chromium
- 静态资源由前端容器内的 Nginx 提供
- 公共上传目录：`/public/uploads`

### 界面与入口
- 模板、编辑器、管理、打印等页面已实现，可直接在应用中体验。
- 打印渲染入口：`frontend/pages/print/PrintResume.tsx:7`

### 快速开始
- 环境要求：
  - Go 1.24+
  - Node.js 20+
  - MySQL 8.x（或 SQLite）
  - Redis 7.x
  - 可用的 Chromium DevTools（推荐 Browserless）

- 本地启动（后端）：
  - 配置环境变量（见 配置）。
  - 启动：`go run ./main.go`
  - 默认端口：`:8080`

- 初始化默认种子数据（岗位分类/岗位/内容预设）：
  - 一次性导入：`go run ./main.go seed import-default`
  - 导入自定义文件：`go run ./main.go seed import --file path/to/seed.json`
  - 从目录导入：`go run ./main.go seed import --dir internal/module/seed/default`
  - 默认数据源文件：`internal/module/seed/default/{categories,roles,presets}.json`
  - `external_id` 是幂等键；发布后尽量保持稳定不变。

- 本地启动（前端）：
  - `cd frontend`
  - `npm install`
  - `npm run dev`
  - Dev 服务器默认地址 `http://localhost:3000`
  - `frontend/config.ts` 读取 `VITE_API_BASE`（默认 `/api/v1`），需正确指向或代理到后端。
  - 开发代理将 `/api` 与 `/public` 转发到 `http://localhost:8080`（见 `frontend/vite.config.ts`）。

- Docker Compose：
  - 根据需要调整 `docker-compose.yml` 的环境变量。
  - 启动：`docker compose up -d`
  - 服务包含：backend、frontend、mysql、redis、chrome。
  - 前端 Nginx 将 `/api/` 与 `/public/` 代理到后端，详见 `frontend/nginx.conf`。

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
- `CHROME_API_URL`：Browserless/Chromium 服务的基础地址，需提供 `/pdf` 与 `/screenshot`
- `enable_email_verification`：注册是否启用邮箱验证码（默认 `false`）
- SMTP：`SMTP_HOST`、`SMTP_PORT`、`SMTP_USERNAME`、`SMTP_PASSWORD`、`SMTP_FROM_NAME`
- 微信 OAuth：`WECHAT_APP_ID`、`WECHAT_APP_SECRET`、`WECHAT_REDIRECT_URI`，开关 `enabled_wechat_login`
- GitHub OAuth：`GITHUB_CLIENT_ID`、`GITHUB_CLIENT_SECRET`、`GITHUB_REDIRECT_URI`，开关 `enabled_github_login`
- S3 存储：`S3_BUCKET`、`S3_REGION`、`S3_ENDPOINT`、`S3_ACCESS_KEY`、`S3_SECRET_KEY`，开关 `enabled_storage_s3`

前端环境变量：
- `VITE_API_BASE`：API 基路径或 URL，默认 `/api/v1`
- `VITE_OAUTH_ALLOWED_ORIGINS`：允许进行 OAuth 弹窗消息通信的来源（逗号分隔）

配置示例：
- 本地后端（bash）：

```bash
export PORT=8080
export DB_DSN=""                            # 为空时使用 SQLite
export SQLITE_PATH="openresume.db"          # SQLite 文件路径
export REDIS_ADDR="127.0.0.1:6379"
export REDIS_PASSWORD=""
export JWT_SECRET="请替换为强随机密钥"
# 首次运行用于注入系统参数（可在后台后续修改）
export CORS_ORIGINS="http://localhost:3000"
export FRONTEND_BASE_URL="http://localhost:3000"
export CHROME_API_URL="http://localhost:3000"  # Browserless/Chromium API 基址，提供 /pdf 与 /screenshot
export SMTP_HOST=""
export SMTP_PORT=""
export SMTP_USERNAME=""
export SMTP_PASSWORD=""
export SMTP_FROM_NAME=""
export WECHAT_APP_ID=""
export WECHAT_APP_SECRET=""
export WECHAT_REDIRECT_URI=""
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
VITE_API_BASE=/api/v1
VITE_OAUTH_ALLOWED_ORIGINS=http://localhost:3000
```

- Docker Compose：
  - 参考 `docker-compose.yml` 中的示例值（已为后端与前端构建参数接好）。
  - 将 `CORS_ORIGINS`、`FRONTEND_BASE_URL`、`CHROME_API_URL` 与前端的 `VITE_OAUTH_ALLOWED_ORIGINS` 替换为你的域名。

### 日志与监控
- 请求日志：使用 Zap 的生产 JSON 格式记录 HTTP 访问，包含 method、path、status、duration、request_id。参考中间件 [logger.go](file:///Users/jaylen/go/src/OpenResume/internal/middleware/logger.go#L11-L17) 与日志封装 [logger.go](file:///Users/jaylen/go/src/OpenResume/internal/pkg/logger/logger.go#L10-L34)。
- 请求 ID 关联：通过 [requestid.go](file:///Users/jaylen/go/src/OpenResume/internal/middleware/requestid.go#L8-L16) 生成/透传 X-Request-ID，并在日志中附带，便于端到端关联。
- 指标：内置 Prometheus 指标，计数 http_requests_total 与直方图 http_request_duration_seconds；端点为 `GET /api/v1/metrics`，见 [metrics.go](file:///Users/jaylen/go/src/OpenResume/internal/pkg/metrics/metrics.go#L12-L29)。
- Loki/Grafana（可选）：提供示例栈用于聚合容器日志与可视化，见 [docker-compose.yaml](file:///Users/jaylen/go/src/OpenResume/loki-stack/docker-compose.yaml) 与 [alloy.river](file:///Users/jaylen/go/src/OpenResume/loki-stack/alloy.river)。

### 中间件
- RequestID：生成/透传 `X-Request-ID`，见 [requestid.go](file:///Users/jaylen/go/src/OpenResume/internal/middleware/requestid.go#L8-L16)
- Logger：HTTP 访问日志，见 [logger.go](file:///Users/jaylen/go/src/OpenResume/internal/middleware/logger.go#L11-L17)
- Recovery：Gin 内置异常恢复
- CORS：按系统配置允许来源，见 [cors.go](file:///Users/jaylen/go/src/OpenResume/internal/middleware/cors.go#L12-L37)
- RateLimit（IP 级）：基于 Redis 的限流，见 [ratelimit.go](file:///Users/jaylen/go/src/OpenResume/internal/middleware/ratelimit.go#L14-L28)
- RateLimitUser（用户级）：登录用户/匿名身份的限流，见 [ratelimit_user.go](file:///Users/jaylen/go/src/OpenResume/internal/middleware/ratelimit_user.go#L15-L34)
- DailyUV：按日去重的 UV 统计，见 [uv.go](file:///Users/jaylen/go/src/OpenResume/internal/middleware/uv.go#L20-L41)
- Auth / RequireRole：JWT 校验与角色控制，见 [auth.go](file:///Users/jaylen/go/src/OpenResume/internal/middleware/auth.go#L15-L54) 与 [admin.go](file:///Users/jaylen/go/src/OpenResume/internal/middleware/admin.go#L13-L41)

### API 概览
- 公共接口：
  - `GET /api/v1/templates`：模板列表
  - `GET /api/v1/templates/:id`：模板详情
  - `GET /api/v1/public/resumes/:slug`：查看公开简历
  - `GET /api/v1/healthz`：健康检查
  - `GET /api/v1/metrics`：指标
  - `GET /api/v1/pdf/exports/:job_id/download`：下载导出结果
- 认证相关：
  - `POST /api/v1/auth/send-code`：发送邮箱验证码
  - `POST /api/v1/auth/register`：邮箱+验证码+密码注册
  - `POST /api/v1/auth/login`：邮箱+密码登录
  - `POST /api/v1/auth/refresh`：刷新令牌
  - `POST /api/v1/auth/logout`：退出登录
  - `GET /api/v1/auth/github/redirect` / `GET /api/v1/auth/github/callback`
  - `GET /api/v1/auth/wechat/redirect` / `GET /api/v1/auth/wechat/callback`
  - `POST /api/v1/auth/wechat/consume-ott`：消费弹窗/重定向的一次性令牌
- 需登录：
  - `GET/POST/PUT/DELETE /api/v1/resumes` 及 `/:id`
  - `POST /api/v1/resumes/:id/image`：生成 PNG
  - `POST /api/v1/resumes/:id/publish`：创建分享链接
  - `POST /api/v1/pdf/exports`：提交导出任务
  - `GET /api/v1/pdf/exports/:job_id`：查询任务状态
  - `GET /api/v1/users/me`、`PUT /api/v1/users/profile`、`PUT /api/v1/users/password`
- 管理接口：
  - 用户、简历、模板、分享链接管理位于 `/api/v1/admin/...`
  - 系统参数可在 `/#/admin/configs` 页面维护

### 导出服务
- 通过 DevTools WebSocket 控制远程 Chromium。
- `CHROME_API_URL` 需指向提供 REST 接口的 Browserless/Chromium 服务：
  - `POST /pdf` 用于生成 PDF
  - `POST /screenshot` 用于生成图片
- 服务会导航到打印路由并执行 PDF 或截图导出；见 `internal/module/pdf/service.go:92`。
- 确保 Chromium 已安装中文字体以正确显示 CJK 内容。
- 提供异步导出队列：提交、查询与下载接口见上方 API；任务由内置工作线程处理与上传，参考 [worker.go](file:///Users/jaylen/go/src/OpenResume/internal/module/pdf/worker.go#L35-L76) 与 [worker.go:StartWorker](file:///Users/jaylen/go/src/OpenResume/internal/module/pdf/worker.go#L89-L104)。

### 国际化
- 内置语言切换，语言偏好持久化在 `localStorage`（键名 `lang`）。
- 文案定义位置：`frontend/contexts/LanguageContext.tsx:10`。
- 默认语言为英文，可在页脚或设置中切换。

### 开发提示
- 后端：
  - 如存在 `db/migrations`，优先使用版本化迁移；否则自动迁移；见 `internal/infra/db/db.go:35`。
  - 断路器状态保存在 Redis 键 `cb:*`。
  - CORS 由 `CORS_ORIGINS` 控制；见 `internal/middleware/cors.go:11`。
- 前端：
  - 路由使用 HashRouter；见 `frontend/App.tsx`。
  - 保持 `VITE_API_BASE` 与后端代理路径一致。
  - 首个注册用户会被赋予 `admin` 角色；之后为 `user`。

### 部署建议
- 推荐使用 `docker-compose`，包含 Nginx 与 Browserless Chromium。
- 通过安全方式管理敏感变量，避免硬编码。
- 生产环境设置 `FRONTEND_BASE_URL` 为你的域名，配置好 `CORS_ORIGINS`，并将前端的 `VITE_OAUTH_ALLOWED_ORIGINS` 传入构建。

### 安全
- 妥善保管 `JWT_SECRET`、SMTP 与 OAuth 密钥。
- 严格限定 `OAUTH_ALLOWED_ORIGINS`。
- 不要将 DevTools 对外公开，建议仅内部访问。

### 参与贡献
- Fork 仓库，创建特性分支，提交 PR 并清晰描述改动。
- 保持改动内聚，必要时更新文档。
- 欢迎提 Issue 与建议。

### 许可证
- MIT 许可。参见 `LICENSE`。
