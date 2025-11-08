# Backend API (Go + Echo)

The Golang service exposes secured REST endpoints for dashboards, orchestrates forecasting/prescriptive calls to the AI microservice, and persists analytics artifacts.

---

## 🔑 Capabilities

- JWT authentication with role/permission middleware.
- Aggregated analytics for Sales, Production, Finance, Inventory, HR, Supply Chain, Executive BI.
- Predictive & prescriptive endpoints (`/api/ai/predict`, `/api/ai/prescribe`, `/api/ai/anomaly`, `/api/ai/scenario`).
- Integration/ETL triggers, business calculators, report exports.
- Request caching helpers and security middleware (headers, rate limiting).

---

## ⚙️ Setup

```bash
cd backend-api
go mod download
cp .env.example .env   # configure DB + AI_SERVICE_URL

go run main.go
```

By default the API listens on `http://localhost:8080`. Confirm health via `GET /api/health`.

---

## 📚 Key Routes

| Group                | Examples |
|----------------------|----------|
| Auth                 | `POST /api/auth/login`, `GET /api/auth/user` |
| Sales                | `GET /api/sales/overview`, `GET /api/sales/channelwise` |
| Finance              | `GET /api/finance/overview`, `GET /api/finance/budget` |
| Inventory            | `GET /api/inventory/overview`, `GET /api/inventory/ratio` |
| Executive BI         | `GET /api/dashboard/executive` |
| AI Predictive        | `GET /api/ai/predict/sales`, `GET /api/ai/prescribe/inventory`, `GET /api/ai/prescribe/finance`, `GET /api/ai/anomaly`, `POST /api/ai/scenario` |
| Integrations         | `POST /api/integration/sync`, `GET /api/integration/status` |

All protected routes require `Authorization: Bearer <token>` header.

---

## 🧭 Development Workflow

1. Register routes in `routes/routes.go` (grouped per domain).
2. Implement controller functions in `controllers/` using `database.DB` (GORM) and AI helpers (`postToAIService`, `callAISummary`).
3. Update models/DTOs under `models/` if schema changes.
4. Run `go fmt ./...` and `go test ./...` before committing.

Detailed patterns are documented in the [Backend Developer Guide](../docs/developer/BACKEND_API_GUIDE.md).

---

## 🧪 Tooling

| Task          | Command |
|---------------|---------|
| Run server    | `go run main.go` |
| Tests         | `go test ./...` |
| Build binary  | `go build -o bin/backend-api ./...` |
| Hot reload    | `air` (optional) |

---

## 🗄️ Database Notes

- Uses MySQL 8 with schema defined in `admin-cms/database/migrations`.
- Update Go structs when migrations add/change columns.
- Connection string pulled from config (`config.AppConfig`).

---

## 📎 Resources

- [Setup Guide](../docs/SETUP_GUIDE.md)
- [Feature Catalog](../docs/FEATURE_CATALOG.md)
- [SRE Runbook](../docs/SRE_RUNBOOK.md)

