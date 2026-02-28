<div align="center">
<img width="1200" height="475" alt="GHBanner" src="https://github.com/user-attachments/assets/0aa67016-6eaf-458a-adb2-6e31a0763ed6" />
</div>

# CVForge 前端运行说明

依赖：Node.js 20+

## 本地运行
- 安装依赖：`npm install`
- 在 `.env.local` 配置环境变量：

```bash
VITE_API_BASE=http://localhost:8080
VITE_OAUTH_ALLOWED_ORIGINS=http://localhost:3000
```

- 启动开发：`npm run dev`（默认 `http://localhost:3000`）
- 如需本地对接后端，可在 Vite 开发服务器中配置转发

## 环境变量
- `VITE_API_BASE`：后端基础路径或 URL
- `VITE_OAUTH_ALLOWED_ORIGINS`：允许进行 OAuth 弹窗消息通信的来源（逗号分隔）

## Docker 部署
- 参考根目录的 `docker-compose.yml`
- 通过构建参数 `VITE_OAUTH_ALLOWED_ORIGINS` 注入前端环境变量
