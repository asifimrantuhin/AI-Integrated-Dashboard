# Setup Guide

This guide walks you through preparing the iDash platform for local development, QA, and production deployments.

---

## 1. Prerequisites

| Component       | Version | Notes |
|-----------------|---------|-------|
| Node.js         | 18.x LTS| Used by the React frontend (Vite) |
| npm / pnpm      | 9.x / 8.x | Either toolchain is supported |
| Go              | 1.21+   | Backend API built with Echo + GORM |
| Python          | 3.10+   | FastAPI AI microservice |
| PHP             | 8.1+    | Laravel Admin CMS |
| Composer        | 2.x     | PHP dependency manager |
| MySQL           | 8.0+    | Shared transactional database |
| Redis (opt.)    | 7.x     | Recommended for caching and queues |
| Docker (opt.)   | 24+     | Simplifies orchestration |

> **Tip**: Use a version manager (e.g. `asdf`, `nvm`, `pyenv`, `gvm`) to align with project defaults.

---

## 2. Environment Variables

Create `.env` files per service by copying the sample template.

```bash
# Frontend
cp frontend/.env.example frontend/.env

# Backend API
tcp backend-api/.env.example backend-api/.env

# Admin CMS
cp admin-cms/.env.example admin-cms/.env

# AI Service
cp ai-service/.env.example ai-service/.env
```

Populate the following minimum keys:

```ini
# Shared
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=idash
DB_USER=idash_user
DB_PASSWORD=supersecret

# Backend API
aI_SERVICE_URL=http://localhost:8080
JWT_SECRET=replace_me
LOG_LEVEL=info

# Admin CMS
APP_URL=http://localhost:8001
QUEUE_CONNECTION=database

# AI Service
AI_LOG_LEVEL=info
```

For production, store secrets in a managed vault (AWS Secrets Manager, HashiCorp Vault, Azure Key Vault) and inject via CI/CD.

---

## 3. Database Preparation

1. **Create the schema**
   ```sql
   CREATE DATABASE idash CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
   CREATE USER 'idash_user'@'%' IDENTIFIED BY 'supersecret';
   GRANT ALL PRIVILEGES ON idash.* TO 'idash_user'@'%';
   FLUSH PRIVILEGES;
   ```
2. **Run Laravel migrations** (covers shared schema for all modules):
   ```bash
   cd admin-cms
   php artisan migrate
   ```
3. **Seed RBAC data**:
   ```bash
   php artisan db:seed --class=RolePermissionSeeder
   ```
4. **Generate demo data** (optional but recommended):
   ```bash
   php artisan idash:demo-data --years=2 --company=1
   ```
   This command refreshes transactional tables (except users) and inserts 1–2 years of realistic multi-module data.

---

## 4. Local Service Bootstrapping

Run each service in its own terminal or via process manager (tmux, `foreman`, PM2).

### Frontend (Vite + React)
```bash
cd frontend
npm install
npm run dev -- --host
```

*Serves http://localhost:5173 by default.*

### Backend API (Go + Echo)
```bash
cd backend-api
go mod download
go run main.go
```

*Serves http://localhost:8080.*

### Admin CMS (Laravel + Blade)
```bash
cd admin-cms
composer install
php artisan serve --host=0.0.0.0 --port=8001
```

*Serves http://localhost:8001.*

### AI Service (FastAPI)
```bash
cd ai-service
python -m venv venv
source venv/bin/activate  # Windows: venv\Scripts\activate
pip install -r requirements.txt
uvicorn main:app --reload --port 8000
```

*Serves http://localhost:8000.*

---

## 5. Running with Docker

A `docker-compose.yml` (if enabled) simplifies orchestration:

```bash
docker compose up --build
docker compose ps
```

Override configuration via `.env.docker` or environment variables in the compose file.

---

## 6. Verification Checklist

1. Visit the frontend, authenticate using seeded credentials.
2. Confirm dashboards load KPI tiles, charts, and AI recommendations without 500 errors.
3. Hit `GET /api/health` (backend) and `/health` (AI service).
4. In Admin CMS, ensure demo data command executes and sync status updates.

---

## 7. Deployment Considerations

- **Backend/API**: Build static Go binaries, deploy behind a reverse proxy (Nginx/Traefik). Enable HTTPS, configure rate limiting middleware.
- **Frontend**: Produce static bundle with `npm run build` and host via CDN or static hosting (S3 + CloudFront, Netlify, Vercel).
- **Admin CMS**: Deploy on PHP-FPM with Nginx or Apache. Configure queue workers for heavy ETL/data generation jobs.
- **AI Service**: Containerize with Uvicorn + Gunicorn (`tiangolo/uvicorn-gunicorn-fastapi`) and autoscale for compute-intensive forecasts.
- **Database**: Use managed MySQL with automated backups and read replicas. Enable PITR for compliance.
- **Observability**: Export logs to centralized stack (ELK, OpenSearch, Datadog). Configure metrics scraping (Prometheus) and tracing (OpenTelemetry).

---

## 8. Useful Commands Reference

| Use Case                         | Command |
|----------------------------------|---------|
| Clear Laravel caches             | `php artisan optimize:clear` |
| Run Go tests                     | `go test ./...` |
| Run frontend unit tests          | `npm test` |
| Start AI service with hot reload | `uvicorn main:app --reload` |
| Trigger integration sync         | `php artisan idash:sync-external --all` |

---

You're now ready to explore predictive dashboards, prescriptive recommendations, and enterprise automation capabilities within iDash.
