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

### Configuration
Backend process environment (directly loaded at startup):
- `PORT`: HTTP port
- `DB_DSN`: MySQL DSN
- `SQLITE_PATH`: Optional SQLite file path (overrides MySQL when set)
- `REDIS_ADDR`, `REDIS_PASSWORD`: Redis connection
- `JWT_SECRET`: HMAC secret for JWT signing
- `CORS_ORIGINS`: Allowed origins (comma separated)
- `FRONTEND_BASE_URL`: Public URL for the frontend (used by export)
- `CHROME_JSON_URL`: DevTools version endpoint, e.g. `http://chrome:3000/json/version`

System configuration keys (stored in DB; seeded from env on first run; editable in Admin → Settings):
- `enable_email_verification`: Enable email verification during registration (default `false`)
- SMTP: `SMTP_HOST`, `SMTP_PORT`, `SMTP_USERNAME`, `SMTP_PASSWORD`, `SMTP_FROM_NAME`
- WeChat OAuth: `WECHAT_APP_ID`, `WECHAT_APP_SECRET`, `WECHAT_REDIRECT_URI`, toggle `enabled_wechat_login`
- GitHub OAuth: `GITHUB_CLIENT_ID`, `GITHUB_CLIENT_SECRET`, `GITHUB_REDIRECT_URI`, toggle `enabled_github_login`
- S3 storage: `S3_BUCKET`, `S3_REGION`, `S3_ENDPOINT`, `S3_ACCESS_KEY`, `S3_SECRET_KEY`, toggle `enabled_storage_s3`

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
  - System configs editable in UI: `/#/admin/configs`

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
  - HashRouter-based routes; see `frontend/App.tsx`.
  - `VITE_API_BASE` should match backend proxy/path.
  - First registered user is promoted to `admin`; subsequent users are `user`.

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
