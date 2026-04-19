# Hopesy Sub2API

这个仓库是 `Wei-Shaw/sub2api` 的 fork。本文档只保留本 fork 的维护和部署说明。

## 1. 发布本项目镜像

本仓库使用 GitHub Actions 工作流 `.github/workflows/publish-image.yml` 发布自己的 GHCR 镜像。

### 触发方式

- 推送符合 `v*` 规则的 tag 时自动构建镜像
- 或在 GitHub Actions 页面手动触发 `Publish Image`

### 发布命令

```bash
git tag v0.1.0
git push origin v0.1.0
```

### 镜像地址

```text
ghcr.io/hopesy/sub2api:v0.1.0
```

生产环境建议优先使用明确版本号，不要长期依赖 `latest`。

---

## 2. 同步上游更改

第一次配置上游仓库：

```bash
git remote add upstream https://github.com/Wei-Shaw/sub2api.git
git fetch upstream
```

以后同步上游：

```bash
git checkout main
git fetch upstream
git merge upstream/main
git push origin main
```

如果有冲突，先在本地解决，再 push 到自己的 fork。

---

## 3. 本地启动与测试

### 3.1 推荐：Docker Compose 本地整套启动

```bash
cd deploy
docker compose -f docker-compose.dev.yml up --build
```

默认访问地址：

```text
http://127.0.0.1:8080
```

### 3.2 后端测试

```bash
cd backend
go test -tags=unit ./...
go test -tags=integration ./...
```

### 3.3 前端开发

```bash
cd frontend
pnpm install
pnpm dev
```

### 3.4 手动启动（不使用 Docker Compose）

适用于你本机已经准备好 PostgreSQL 和 Redis 的场景。

后端：

```bash
cd backend
go run ./cmd/server
```

前端：

```bash
cd frontend
pnpm install
pnpm dev
```

默认开发态前端地址：

```text
http://127.0.0.1:3000
```

默认后端地址：

```text
http://127.0.0.1:8080
```

如果后端需要读取本地配置，优先准备 `backend/config.yaml`，或按项目约定设置数据库、Redis、JWT 等环境变量。

---

## 4. 部署到 ClawCloud

### 4.1 镜像

在 ClawCloud 中使用本仓库发布的 GHCR 镜像，例如：

```text
ghcr.io/hopesy/sub2api:v0.1.0
```

### 4.2 环境变量

仓库根目录提供了当前可直接参考的环境变量文件：

```text
clawcloud.env
```

导入时重点注意：

- `DATABASE_DBNAME=sub2api`
- `REDIS_HOST` 只能填纯域名，不能带 `https://`
- Upstash 场景需要：
  - `REDIS_PORT=6379`
  - `REDIS_ENABLE_TLS=true`
  - `REDIS_DB=0`

### 4.3 访问地址

ClawCloud 部署成功后，直接访问应用根路径 `/` 即可打开管理后台。

---

## 5. 当前维护重点

本 fork 当前已经针对部署问题做过收口：

- 自动初始化时支持更稳的 PostgreSQL 建库流程
- Redis host 输入兼容处理更稳
- 镜像发布改为 tag 才触发
