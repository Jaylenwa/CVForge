<div align="center">
  <h1>OpenResume</h1>
  <p>AI-powered, bilingual resume builder with templates, exports, sharing, and admin tools.</p>
  <p>
    <a href="README.md">English</a> ·
    <a href="README.zh-CN.md">简体中文</a>
  </p>
</div>

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
 - Logging & monitoring: Zap JSON logs, Prometheus metrics, optional Loki/Grafana.
 
 ### Architecture
  - Frontend: `frontend/` (React 19, Vite 6, TypeScript)
  - Backend: `internal/` modules (Gin handlers, services, middleware), `main.go`
  - Storage: MySQL/SQLite + Redis
- Export: Remote Chromium (Browserless) via WebSocket DevTools
- Static assets served via Nginx in the frontend container
 - Public uploads served from `/public/uploads`

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

- Initialize default seed data (taxonomy + presets):
  - Run once: `go run ./main.go seed import-default`
  - Import from a custom file: `go run ./main.go seed import --file path/to/seed.json`
  - Import from a directory: `go run ./main.go seed import --dir internal/module/seed/default`
  - Default source files: `internal/module/seed/default/{categories,roles,presets}.json`
  - `external_id` is the idempotent key; keep it stable once published.

- Run locally (frontend):
  - `cd frontend`
  - `npm install`
  - `npm run dev`
  - The dev server runs on `http://localhost:3000` by default.
  - `frontend/config.ts` reads `VITE_API_BASE` (default `/api/v1`), which should proxy or point to the backend.
  - Dev proxy forwards `/api` and `/public` to `http://localhost:8080` (see `frontend/vite.config.ts`).

 - Docker Compose:
   - Review and customize `docker-compose.yml` environment variables.
   - Start: `docker compose up -d`
   - Services: backend, frontend, mysql, redis, chrome.
   - Frontend Nginx proxies `/api/` and `/public/` to backend; see `frontend/nginx.conf`.
 
 ### Logging & Monitoring
 - Request logging: Zap production JSON logs for HTTP access, including method, path, status, duration, request_id. See middleware [logger.go](file:///Users/jaylen/go/src/OpenResume/internal/middleware/logger.go#L11-L17) and logger wrapper [logger.go](file:///Users/jaylen/go/src/OpenResume/internal/pkg/logger/logger.go#L10-L34).
 - Request ID propagation: `X-Request-ID` is generated/forwarded by [requestid.go](file:///Users/jaylen/go/src/OpenResume/internal/middleware/requestid.go#L8-L16) and attached to logs for end-to-end tracing.
 - Metrics: Built-in Prometheus metrics with `http_requests_total` and `http_request_duration_seconds`; endpoint `GET /api/v1/metrics`. See [metrics.go](file:///Users/jaylen/go/src/OpenResume/internal/pkg/metrics/metrics.go#L12-L29).
 - Optional Loki/Grafana stack: Example setup for container log aggregation and dashboards in [docker-compose.yaml](file:///Users/jaylen/go/src/OpenResume/loki-stack/docker-compose.yaml) and [alloy.river](file:///Users/jaylen/go/src/OpenResume/loki-stack/alloy.river).
 
 ### Middlewares
 - RequestID: generate/forward `X-Request-ID`, see [requestid.go](file:///Users/jaylen/go/src/OpenResume/internal/middleware/requestid.go#L8-L16)
 - Logger: HTTP access logging, see [logger.go](file:///Users/jaylen/go/src/OpenResume/internal/middleware/logger.go#L11-L17)
 - Recovery: Gin’s built-in panic recovery
 - CORS: allow origins per system config, see [cors.go](file:///Users/jaylen/go/src/OpenResume/internal/middleware/cors.go#L12-L37)
 - RateLimit (IP-based): Redis-backed limiting, see [ratelimit.go](file:///Users/jaylen/go/src/OpenResume/internal/middleware/ratelimit.go#L14-L28)
 - RateLimitUser (user-based): limits by logged-in user or IP, see [ratelimit_user.go](file:///Users/jaylen/go/src/OpenResume/internal/middleware/ratelimit_user.go#L15-L34)
 - DailyUV: per-day unique visitor counting, see [uv.go](file:///Users/jaylen/go/src/OpenResume/internal/middleware/uv.go#L20-L41)
 - Auth / RequireRole: JWT validation and role gating, see [auth.go](file:///Users/jaylen/go/src/OpenResume/internal/middleware/auth.go#L15-L54) and [admin.go](file:///Users/jaylen/go/src/OpenResume/internal/middleware/admin.go#L13-L41)
 
 ### Configuration
 Backend process environment (read at startup):
 - `PORT`: HTTP port
 - `DB_DSN`: MySQL DSN; when set, MySQL is used
- `SQLITE_PATH`: SQLite file path; used when `DB_DSN` is empty
- `REDIS_ADDR`, `REDIS_PASSWORD`: Redis connection
- `JWT_SECRET`: JWT signing secret

System configuration keys (stored in DB; seeded from env on first run; editable in Admin → Settings):
- `CORS_ORIGINS`: Allowed origins (comma separated)
- `FRONTEND_BASE_URL`: Public frontend URL (used by export and OAuth redirects)
- `CHROME_API_URL`: Base URL of your Browserless/Chromium service providing `/pdf` and `/screenshot`
- `enable_email_verification`: Enable email verification during registration (default `false`)
- SMTP: `SMTP_HOST`, `SMTP_PORT`, `SMTP_USERNAME`, `SMTP_PASSWORD`, `SMTP_FROM_NAME`
- WeChat OAuth: `WECHAT_APP_ID`, `WECHAT_APP_SECRET`, `WECHAT_REDIRECT_URI`, toggle `enabled_wechat_login`
- WeChat MP QR login (follow to sign in): `WECHAT_MP_APP_ID`, `WECHAT_MP_APP_SECRET`, `WECHAT_MP_TOKEN` (optional `WECHAT_MP_AES_KEY`), toggle `enabled_wechat_mp_login`; callback URL is `GET/POST /api/v1/wechat/mp/callback` (must be publicly reachable over HTTPS)
- GitHub OAuth: `GITHUB_CLIENT_ID`, `GITHUB_CLIENT_SECRET`, `GITHUB_REDIRECT_URI`, toggle `enabled_github_login`
- S3 storage: `S3_BUCKET`, `S3_REGION`, `S3_ENDPOINT`, `S3_ACCESS_KEY`, `S3_SECRET_KEY`, toggle `enabled_storage_s3`

Environment variables (frontend):
- `VITE_API_BASE`: Base path or URL for API, default `/api/v1`
- `VITE_OAUTH_ALLOWED_ORIGINS`: Comma-separated origins allowed for OAuth popup messaging

Configuration examples:
- Local backend (bash):

```bash
export PORT=8080
export DB_DSN=""                            # use SQLite when empty
export SQLITE_PATH="openresume.db"          # SQLite file path
export REDIS_ADDR="127.0.0.1:6379"
export REDIS_PASSWORD=""
export JWT_SECRET="replace-with-a-strong-secret"
# Seed system configs on first run (optional; can be edited later in Admin → Settings)
export CORS_ORIGINS="http://localhost:3000"
export FRONTEND_BASE_URL="http://localhost:3000"
export CHROME_API_URL="http://localhost:3000"  # Browserless/Chromium API base providing /pdf and /screenshot
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

- Frontend `.env.local`:

```bash
VITE_API_BASE=/api/v1
VITE_OAUTH_ALLOWED_ORIGINS=http://localhost:3000
```

- Docker Compose:
  - See `docker-compose.yml` for example values already wired for backend and frontend build args.
  - Replace `CORS_ORIGINS`, `FRONTEND_BASE_URL`, `CHROME_API_URL` and `VITE_OAUTH_ALLOWED_ORIGINS` with your domain(s).

 ### API Overview
 - Public:
   - `GET /api/v1/templates` — list templates
   - `GET /api/v1/templates/:id` — template details
   - `GET /api/v1/public/resumes/:slug` — view public resume
   - `GET /api/v1/healthz` — health
   - `GET /api/v1/metrics` — metrics
   - `GET /api/v1/pdf/exports/:job_id/download` — download export result
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
   - `POST /api/v1/resumes/:id/image` — generate PNG
   - `POST /api/v1/resumes/:id/publish` — create public share link
   - `POST /api/v1/pdf/exports` — submit export job
   - `GET /api/v1/pdf/exports/:job_id` — check export job status
   - `GET /api/v1/users/me`, `PUT /api/v1/users/profile`, `PUT /api/v1/users/password`
 - Admin:
   - Users, Resumes, Templates, Shares management under `/api/v1/admin/...`
   - System configs editable in UI: `/#/admin/configs`
 
 ### Export Services
 - Uses a remote Chromium via DevTools WebSocket.
 - `CHROME_API_URL` must point to a Browserless/Chromium service that exposes REST endpoints:
   - `POST /pdf` for PDF generation
   - `POST /screenshot` for image generation
 - The service navigates to the Print route and executes `page.PrintToPDF` or screenshot; see `internal/module/pdf/service.go:92`.
 - Ensure Chinese fonts are installed in Chromium for CJK content.
 - Asynchronous export queue: submit, query, and download via the endpoints above. Jobs are processed by an in-process worker and uploaded to storage; see [worker.go](file:///Users/jaylen/go/src/OpenResume/internal/module/pdf/worker.go#L35-L76) and [worker.go:StartWorker](file:///Users/jaylen/go/src/OpenResume/internal/module/pdf/worker.go#L89-L104).

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
  - HashRouter-based routes; see `frontend/App.tsx`.
  - `VITE_API_BASE` should match backend proxy/path.
  - First registered user is promoted to `admin`; subsequent users are `user`.

### Deployment
- Recommended: `docker-compose` with Nginx and Browserless Chromium.
- Use secrets management for sensitive variables (do not hardcode).
- For production, set `FRONTEND_BASE_URL` to your domain, configure `CORS_ORIGINS`, and pass `VITE_OAUTH_ALLOWED_ORIGINS` to the frontend.

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
