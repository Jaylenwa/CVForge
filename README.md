<div align="center">
  <h1>OpenResume</h1>
  <p>AI-powered, bilingual resume builder with templates, exports, sharing, and admin tools.</p>
  <p>
    <a href="#english">English</a> ·
    <a href="#简体中文">简体中文</a>
  </p>
</div>

## English

### What Is OpenResume
- Full-stack resume builder focused on speed, ATS-friendly design, and strong UX.
- AI assists with writing and polishing content.
- Export to PDF/PNG using a remote Chromium instance.
- Share public resume links; manage visibility and analytics.
- Admin panel for users, resumes, templates, and share links.
- Built-in internationalization with English and Chinese.

### Highlights
- Modern frontend with React + Vite + TypeScript.
- Go backend using Gin, Gorm, Redis, JWT.
- OAuth login for GitHub and WeChat; email code verification via SMTP.
- Docker-first deployments; `docker-compose` for production.
- Circuit breaker around export services for stability.

### Architecture
- Frontend: `frontend/` (React 19, Vite 6, TypeScript)
- Backend: `internal/` modules (Gin handlers, services, middleware), `main.go`
- Storage: MySQL/SQLite + Redis
- Export: Remote Chromium (Browserless) via WebSocket DevTools
- Static assets served via Nginx in the frontend container

Key router entry: `internal/router/router.go:30` (`/api/v1` routes)  
Config loading: `internal/infra/config/config.go:40`  
PDF/PNG generation: `internal/module/pdf/service.go:57`  
AI endpoints: `internal/module/ai/handler.go:26`  
Auth endpoints: `internal/module/auth/handler.go:229`

### Screenshots
- Templates, Editor, Admin, and Print pages are implemented; you can explore via the app’s UI.
- Print rendering entry: `frontend/pages/print/PrintResume.tsx:7`

### Quick Start
- Requirements:
  - Go 1.24+
  - Node.js 20+
  - MySQL 8.x (or SQLite)
  - Redis 7.x
  - A Chromium instance with DevTools (Browserless recommended)

- Run locally (backend):
  - Set environment variables (see Configuration).
  - Start: `go run ./main.go`
  - Default port: `:8080`

- Run locally (frontend):
  - `cd frontend`
  - `npm install`
  - `npm run dev`
  - The dev server runs on `http://localhost:5173` by default.
  - `frontend/config.ts` reads `VITE_API_BASE` (default `/api/v1`), which should proxy or point to the backend.

- Docker Compose:
  - Review and customize `docker-compose.yml` environment variables.
  - Start: `docker compose up -d`
  - Services: backend, frontend, mysql, redis, chrome.
  - Frontend Nginx proxies `/api/` and `/public/` to backend; see `frontend/nginx.conf`.

### Configuration
Environment variables (backend) from `internal/infra/config/config.go`:
- `PORT`: HTTP port
- `DB_DSN`: MySQL DSN
- `SQLITE_PATH`: Optional SQLite file path (overrides MySQL when set)
- `REDIS_ADDR`, `REDIS_PASSWORD`: Redis connection
- `JWT_SECRET`: HMAC secret for JWT signing
- `CORS_ORIGINS`: Allowed origins (comma separated)
- `GEMINI_API_KEY`: Optional key for AI features (if integrated)
- `UPLOAD_BACKEND`: `local` or `s3`
- `S3_BUCKET`, `S3_REGION`, `S3_ENDPOINT`, `S3_ACCESS_KEY`, `S3_SECRET_KEY`: S3 config
- `FRONTEND_BASE_URL`: Public URL for the frontend (used by export)
- `CHROME_JSON_URL`: DevTools version endpoint, e.g. `http://chrome:3000/json/version`
- `SMTP_HOST`, `SMTP_PORT`, `SMTP_USERNAME`, `SMTP_PASSWORD`, `SMTP_FROM_NAME`: SMTP config for email verification
- `WECHAT_APP_ID`, `WECHAT_APP_SECRET`, `WECHAT_REDIRECT_URI`: WeChat OAuth
- `GITHUB_CLIENT_ID`, `GITHUB_CLIENT_SECRET`, `GITHUB_REDIRECT_URI`: GitHub OAuth
- `OAUTH_ALLOWED_ORIGINS`: Allowed origins for OAuth popups/redirects
- `FEATURE_WECHAT_LOGIN`, `FEATURE_GITHUB_LOGIN`: `on`/`off` feature switches

Environment variables (frontend):
- `VITE_API_BASE`: Base path or URL for API, default `/api/v1`
- `VITE_OAUTH_ALLOWED_ORIGINS`: Comma-separated origins used in OAuth flows

### API Overview
- Public:
  - `GET /api/v1/templates` — list templates
  - `GET /api/v1/public/resumes/:slug` — view public resume
  - `GET /api/v1/healthz` — health
  - `GET /api/v1/metrics` — metrics
- Auth:
  - `POST /api/v1/auth/send-code` — email verification code
  - `POST /api/v1/auth/register` — register with email/code/password
  - `POST /api/v1/auth/login` — email/password login
  - `POST /api/v1/auth/refresh` — refresh tokens
  - `POST /api/v1/auth/logout` — logout
  - `GET /api/v1/auth/github/redirect` / `GET /api/v1/auth/github/callback`
  - `GET /api/v1/auth/wechat/redirect` / `GET /api/v1/auth/wechat/callback`
  - `POST /api/v1/auth/wechat/consume-ott` — consume one-time token from popup/redirect flow
- Authenticated:
  - `GET/POST/PUT/DELETE /api/v1/resumes` and `/:id`
  - `POST /api/v1/resumes/:id/pdf` — generate PDF
  - `POST /api/v1/resumes/:id/image` — generate PNG
  - `POST /api/v1/resumes/:id/publish` — create public share link
  - `GET /api/v1/users/me`, `PUT /api/v1/users/profile`, `PUT /api/v1/users/password`
- Admin:
  - Users, Resumes, Templates, Shares management under `/api/v1/admin/...`

### Export Services
- Uses a remote Chromium via DevTools WebSocket.
- `CHROME_JSON_URL` must point to a `json/version` endpoint that returns `webSocketDebuggerUrl`.
- The service navigates to the Print route and executes `page.PrintToPDF` or screenshot; see `internal/module/pdf/service.go:92`.
- Ensure Chinese fonts are installed in Chromium for CJK content.

### Internationalization
- Built-in language switching persists in `localStorage` (`lang` key).
- UI strings in `frontend/contexts/LanguageContext.tsx:10`.
- Default language is English; toggle in the UI footer or settings.

### Development
- Backend:
  - Database initialization uses AutoMigrate unless `db/migrations` exist; see `internal/infra/db/db.go:35`.
  - Circuit breaker state is stored in Redis keys `cb:*`.
  - CORS is configured per `CORS_ORIGINS`; see `internal/middleware/cors.go:11`.
- Frontend:
  - HashRouter-based routes; see `frontend/App.tsx:34`.
  - `VITE_API_BASE` should match backend proxy/path.

### Deployment
- Recommended: `docker-compose` with Nginx and Browserless Chromium.
- Use secrets management for sensitive variables (do not hardcode).
- For production, point `FRONTEND_BASE_URL` to your domain and set `OAUTH_ALLOWED_ORIGINS`.

### Security
- Store `JWT_SECRET`, SMTP and OAuth secrets securely.
- Restrict `OAUTH_ALLOWED_ORIGINS` to trusted domains.
- Keep Chromium behind your network; avoid exposing DevTools publicly.

### Contributing
- Fork, create a feature branch, and submit a PR with a clear title and description.
- Keep changes cohesive; update docs where relevant.
- Issues and feature requests are welcome.

### License
- MIT License. See `LICENSE`.

---

## 简体中文

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

关键路由：`internal/router/router.go:30`（`/api/v1` 路由分组）  
配置加载：`internal/infra/config/config.go:40`  
PDF/PNG 生成：`internal/module/pdf/service.go:57`  
AI 接口：`internal/module/ai/handler.go:26`  
认证接口：`internal/module/auth/handler.go:229`

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
  - Dev 服务器默认地址 `http://localhost:5173`
  - `frontend/config.ts` 读取 `VITE_API_BASE`（默认 `/api/v1`），需正确指向或代理到后端。

- Docker Compose：
  - 根据需要调整 `docker-compose.yml` 的环境变量。
  - 启动：`docker compose up -d`
  - 服务包含：backend、frontend、mysql、redis、chrome。
  - 前端 Nginx 将 `/api/` 与 `/public/` 代理到后端，详见 `frontend/nginx.conf`。

### 配置
后端环境变量（来自 `internal/infra/config/config.go`）：
- `PORT`：HTTP 端口
- `DB_DSN`：MySQL DSN
- `SQLITE_PATH`：SQLite 文件路径（设置后优先生效）
- `REDIS_ADDR`、`REDIS_PASSWORD`：Redis 连接
- `JWT_SECRET`：JWT HMAC 签名密钥
- `CORS_ORIGINS`：允许的跨域来源（逗号分隔）
- `GEMINI_API_KEY`：可选，用于 AI 能力（如集成）
- `UPLOAD_BACKEND`：`local` 或 `s3`
- `S3_BUCKET`、`S3_REGION`、`S3_ENDPOINT`、`S3_ACCESS_KEY`、`S3_SECRET_KEY`：S3 配置
- `FRONTEND_BASE_URL`：前端的公共地址（导出功能需要）
- `CHROME_JSON_URL`：DevTools 版本接口，例如 `http://chrome:3000/json/version`
- `SMTP_HOST`、`SMTP_PORT`、`SMTP_USERNAME`、`SMTP_PASSWORD`、`SMTP_FROM_NAME`：SMTP 邮件配置
- `WECHAT_APP_ID`、`WECHAT_APP_SECRET`、`WECHAT_REDIRECT_URI`：微信 OAuth
- `GITHUB_CLIENT_ID`、`GITHUB_CLIENT_SECRET`、`GITHUB_REDIRECT_URI`：GitHub OAuth
- `OAUTH_ALLOWED_ORIGINS`：OAuth 弹窗/跳转的允许来源
- `FEATURE_WECHAT_LOGIN`、`FEATURE_GITHUB_LOGIN`：`on`/`off` 功能开关

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
  - 路由使用 HashRouter；见 `frontend/App.tsx:34`。
  - 保持 `VITE_API_BASE` 与后端代理路径一致。

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

