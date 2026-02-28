<div align="center">
  <h1>CVForge</h1>
  <p>AI-powered, bilingual resume builder with templates, exports, sharing, and admin tools.</p>
  <p>
    <a href="README.md">English</a> ·
    <a href="README.zh-CN.md">简体中文</a>
  </p>
</div>

### What Is CVForge
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
 - Permissions: first registered user becomes `admin`; subsequent users are `user`.
 - Docker-first deployments; `docker-compose` for production.
 - Circuit breaker around export services for stability.
 - Logging & monitoring: Zap JSON logs, Prometheus metrics, optional Loki/Grafana.
 
 ### Architecture
  - Frontend: React + Vite + TypeScript
  - Backend: Go + Gin (auth, admin, services)
  - Storage: MySQL/SQLite + Redis (optional S3-compatible storage)
  - Export: Remote Chromium (Browserless recommended)
  - Static assets: Nginx

### Screenshots
- Templates, Editor, Admin, and Print pages are implemented; you can explore via the app’s UI.

### Quick Start
- Requirements:
  - Go 1.24+
  - Node.js 20+
  - MySQL 8.x (or SQLite)
  - Redis 7.x
  - A Chromium instance with DevTools (Browserless recommended)

- Run locally (backend):
  - Set environment variables (see Configuration).
  - Start: `go run .`
  - Default port: `:8080`

- Initialize default seed data (taxonomy + presets):
  - Run once: `go run . seed import-default`
  - For custom seed data, use the `seed` subcommand to import from a file.

- Run locally (frontend):
  - `cd frontend`
  - `npm install`
  - `npm run dev`
  - The dev server runs on `http://localhost:3000` by default.
  - Configure the backend base URL and OAuth messaging origins via `.env.local`.

 - Docker Compose:
   - Review and customize `docker-compose.yml` environment variables.
   - Start: `docker compose up -d`
   - Services: backend, frontend, mysql, redis, chrome.
   - Frontend Nginx is preconfigured to proxy backend routes and serve public assets.
 
 ### Logging & Monitoring
 - Request logging: Zap JSON logs for HTTP access (method, path, status, duration, request_id).
 - Request tracing: request ID propagation for end-to-end debugging.
 - Metrics: built-in Prometheus metrics (requests count + duration histogram).
 - Optional Loki/Grafana stack: example setup for container log aggregation and dashboards.
 
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
- `CHROME_API_URL`: Base URL of your Browserless/Chromium service (used for PDF/PNG export)
- `enable_email_verification`: Enable email verification during registration (default `false`)
- SMTP: `SMTP_HOST`, `SMTP_PORT`, `SMTP_USERNAME`, `SMTP_PASSWORD`, `SMTP_FROM_NAME`
- WeChat OAuth: `WECHAT_APP_ID`, `WECHAT_APP_SECRET`, `WECHAT_REDIRECT_URI`, toggle `enabled_wechat_login`
- WeChat MP QR login (follow to sign in): `WECHAT_MP_APP_ID`, `WECHAT_MP_APP_SECRET`, `WECHAT_MP_TOKEN` (optional `WECHAT_MP_AES_KEY`), toggle `enabled_wechat_mp_login`; the callback route must be publicly reachable over HTTPS
- GitHub OAuth: `GITHUB_CLIENT_ID`, `GITHUB_CLIENT_SECRET`, `GITHUB_REDIRECT_URI`, toggle `enabled_github_login`
- S3 storage: `S3_BUCKET`, `S3_REGION`, `S3_ENDPOINT`, `S3_ACCESS_KEY`, `S3_SECRET_KEY`, toggle `enabled_storage_s3`

Environment variables (frontend):
- `VITE_API_BASE`: Backend base path or URL
- `VITE_OAUTH_ALLOWED_ORIGINS`: Comma-separated origins allowed for OAuth popup messaging

Configuration examples:
- Local backend (bash):

```bash
export PORT=8080
export DB_DSN=""                            # use SQLite when empty
export SQLITE_PATH="cvforge.db"          # SQLite file path
export REDIS_ADDR="127.0.0.1:6379"
export REDIS_PASSWORD=""
export JWT_SECRET="replace-with-a-strong-secret"
# Seed system configs on first run (optional; can be edited later in Admin → Settings)
export CORS_ORIGINS="http://localhost:3000"
export FRONTEND_BASE_URL="http://localhost:3000"
export CHROME_API_URL="http://localhost:3000"  # Browserless/Chromium base for PDF/PNG export
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
VITE_API_BASE=http://localhost:8080
VITE_OAUTH_ALLOWED_ORIGINS=http://localhost:3000
```

- Docker Compose:
  - See `docker-compose.yml` for example values already wired for backend and frontend build args.
  - Replace `CORS_ORIGINS`, `FRONTEND_BASE_URL`, `CHROME_API_URL` and `VITE_OAUTH_ALLOWED_ORIGINS` with your domain(s).
 
 ### Export Services
 - Uses a remote Chromium via DevTools WebSocket.
 - `CHROME_API_URL` must point to a Browserless/Chromium service used for PDF/PNG generation.
 - The service opens a print page and executes PDF export or screenshot generation.
 - Ensure Chinese fonts are installed in Chromium for CJK content.
 - Asynchronous export queue: jobs are processed by an in-process worker and uploaded to storage.

### Internationalization
- Built-in language switching (English / Chinese).
- Default language is English; toggle in the UI footer or settings.

### Deployment
- Recommended: `docker-compose` with Nginx and Browserless Chromium.
- Use secrets management for sensitive variables (do not hardcode).
- For production, set `FRONTEND_BASE_URL` to your domain, configure `CORS_ORIGINS`, and pass `VITE_OAUTH_ALLOWED_ORIGINS` to the frontend.

### Security
- Store `JWT_SECRET`, SMTP and OAuth secrets securely.
- Restrict `VITE_OAUTH_ALLOWED_ORIGINS` to trusted domains.
- Keep Chromium behind your network; avoid exposing DevTools publicly.

### Contributing
- Fork, create a feature branch, and submit a PR with a clear title and description.
- Keep changes cohesive; update docs where relevant.
- Issues and feature requests are welcome.

### License
- MIT License. See `LICENSE`.
