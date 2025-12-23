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
后端启动时直接读取的环境变量：
- `PORT`：HTTP 端口
- `DB_DSN`：MySQL DSN
- `SQLITE_PATH`：SQLite 文件路径（设置后优先生效）
- `REDIS_ADDR`、`REDIS_PASSWORD`：Redis 连接
- `JWT_SECRET`：JWT HMAC 签名密钥
- `CORS_ORIGINS`：允许的跨域来源（使用逗号分隔）
- `FRONTEND_BASE_URL`：前端的公共地址（导出功能需要）
- `CHROME_JSON_URL`：DevTools 版本接口，例如 `http://chrome:3000/json/version`

系统配置键（存储于数据库；首次启动从环境变量注入默认值；可在管理后台修改）：
- `enable_email_verification`：注册时是否启用邮箱验证码（默认 `false`）
- SMTP：`SMTP_HOST`、`SMTP_PORT`、`SMTP_USERNAME`、`SMTP_PASSWORD`、`SMTP_FROM_NAME`
- 微信 OAuth：`WECHAT_APP_ID`、`WECHAT_APP_SECRET`、`WECHAT_REDIRECT_URI`，开关 `wechat_login`
- GitHub OAuth：`GITHUB_CLIENT_ID`、`GITHUB_CLIENT_SECRET`、`GITHUB_REDIRECT_URI`，开关 `github_login`
- S3 存储：`S3_BUCKET`、`S3_REGION`、`S3_ENDPOINT`、`S3_ACCESS_KEY`、`S3_SECRET_KEY`，开关 `storage_s3_enabled`

前端环境变量：
- `VITE_API_BASE`：API 基路径或 URL，默认 `/api/v1`
- `VITE_OAUTH_ALLOWED_ORIGINS`：OAuth 流程允许来源列表

### API 概览
- 公共接口：
  - `GET /api/v1/templates`：模板列表
  - `GET /api/v1/public/resumes/:slug`：查看公开简历
  - `GET /api/v1/healthz`：健康检查
  - `GET /api/v1/metrics`：指标
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
  - `POST /api/v1/resumes/:id/pdf`：生成 PDF
  - `POST /api/v1/resumes/:id/image`：生成 PNG
  - `POST /api/v1/resumes/:id/publish`：创建分享链接
  - `GET /api/v1/users/me`、`PUT /api/v1/users/profile`、`PUT /api/v1/users/password`
- 管理接口：
  - 用户、简历、模板、分享链接管理位于 `/api/v1/admin/...`
  - 系统参数可在 `/#/admin/configs` 页面维护

### 导出服务
- 通过 DevTools WebSocket 控制远程 Chromium。
- `CHROME_JSON_URL` 必须返回 `webSocketDebuggerUrl`。
- 服务会导航到打印路由并执行 PDF 或截图导出；见 `internal/module/pdf/service.go:92`。
- 确保 Chromium 已安装中文字体以正确显示 CJK 内容。

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
- 生产环境将 `FRONTEND_BASE_URL` 指向你的域名，并设定 `OAUTH_ALLOWED_ORIGINS`。

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
