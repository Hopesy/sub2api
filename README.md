# Hopesy Sub2API

This repository is a fork of `Wei-Shaw/sub2api`. This document only keeps the operational notes for this fork.

## 1. Publish This Fork Image

This fork publishes its own GHCR image through `.github/workflows/publish-image.yml`.

### Triggers

- Push a tag matching `v*`
- Or run `Publish Image` manually from GitHub Actions

### Release commands

```bash
git tag v0.1.0
git push origin v0.1.0
```

### Image

```text
ghcr.io/hopesy/sub2api:v0.1.0
```

Use an explicit version tag for production instead of relying on `latest` long term.

---

## 2. Sync Upstream Changes

Initial upstream setup:

```bash
git remote add upstream https://github.com/Wei-Shaw/sub2api.git
git fetch upstream
```

Sync later:

```bash
git checkout main
git fetch upstream
git merge upstream/main
git push origin main
```

Resolve conflicts locally before pushing to your fork.

---

## 3. Local Run and Test

### 3.1 Recommended: full local stack with Docker Compose

```bash
cd deploy
docker compose -f docker-compose.dev.yml up --build
```

Default local URL:

```text
http://127.0.0.1:8080
```

### 3.2 Backend tests

```bash
cd backend
go test -tags=unit ./...
go test -tags=integration ./...
```

### 3.3 Frontend development

```bash
cd frontend
pnpm install
pnpm dev
```

### 3.4 Manual startup (without Docker Compose)

Use this when PostgreSQL and Redis are already available on your machine.

Backend:

```bash
cd backend
go run ./cmd/server
```

Frontend:

```bash
cd frontend
pnpm install
pnpm dev
```

Default frontend dev URL:

```text
http://127.0.0.1:3000
```

Default backend URL:

```text
http://127.0.0.1:8080
```

If the backend needs local configuration, prepare `backend/config.yaml` first, or provide the required database, Redis, JWT and related environment variables.

---

## 4. Deploy to ClawCloud

### 4.1 Image

Use your own GHCR image in ClawCloud, for example:

```text
ghcr.io/hopesy/sub2api:v0.1.0
```

### 4.2 Environment variables

The repository root contains a working reference file:

```text
clawcloud.env
```

Important points:

- `DATABASE_DBNAME=sub2api`
- `REDIS_HOST` must be a plain host name, without `https://`
- For Upstash:
  - `REDIS_PORT=6379`
  - `REDIS_ENABLE_TLS=true`
  - `REDIS_DB=0`

### 4.3 Access

After deployment succeeds, open the application root path `/` to access the admin UI.

---

## 5. Current Fork Focus

This fork already contains deployment-oriented fixes for:

- safer PostgreSQL bootstrap when the target database does not exist
- more tolerant Redis host normalization
- tag-only image publishing
