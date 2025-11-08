# Backend API Developer Guide

The Go (Echo) backend exposes REST endpoints consumed by the frontend and admin CMS, and proxies requests to the AI service. This document outlines the development workflow from route registration through response.

---

## 1. Architecture Snapshot

- **Framework**: Echo v4
- **ORM**: GORM
- **Database**: MySQL 8
- **Auth**: JWT via middleware and `auth.go`
- **AI Integration**: HTTP client to FastAPI (`AI_SERVICE_URL`)

Directory structure:
```
backend-api/
├── controllers/    # Route handlers
├── middleware/     # Auth, security headers, rate limiting
├── models/         # GORM models & DTOs
├── routes/         # Route registration
├── services/       # External connectors (integration jobs)
├── utils/          # JWT, caching helpers
└── main.go         # Server entrypoint
```

---

## 2. Request Flow (Route → Response)

1. **Route Definition** (`routes/routes.go`)
   - Grouped by domain (`/sales`, `/inventory`, `/ai`, etc.).
   - Apply `middleware.AuthMiddleware()` globally to require JWT.
   - Use `middleware.RequireRoles(...)` or `RequirePermissions(...)` for finer control.

2. **Controller Logic** (`controllers/*.go`)
   - Fetch query/body parameters via Echo context (`c.QueryParam`, `c.Bind`).
   - Compose DB queries using `database.DB` (GORM).
   - For predictive features, call helper functions such as `callAISummary` or `postToAIService` to reach the AI service, store results, and enrich responses.

3. **Model Mapping** (`models/*.go`)
   - Data retrieved via GORM maps to struct definitions (with JSON tags).
   - Use separate DTO structs when returning aggregated data (e.g., `SalesOverviewResponse`).

4. **Response**
   - Return JSON using `return c.JSON(http.StatusOK, payload)`.
   - Handle errors with consistent structure: `map[string]string{"error": "message"}` and meaningful status codes.

---

## 3. Coding Guidelines

- **Package imports**: follow Go formatting; run `go fmt ./...` before committing.
- **Error handling**: always check returned errors, log with context when necessary.
- **Dependency injection**: rely on global `database.DB` initialization in `main.go` for now; future improvements can use struct-based injection.
- **Transactions**: use `database.DB.Transaction(func(tx *gorm.DB) error { ... })` for multi-step writes.
- **Caching**: `utils/cache.go` provides helpers for storing computed responses (e.g., manufacturing analytics). Use for expensive queries with limited staleness tolerance.
- **AI calls**: `postToAIService` handles request/response marshalling and HTTP status checking.

---

## 4. Adding a New Endpoint (Example)

1. **Define route** in `routes/routes.go`:
   ```go
   api.GET("/sales/performance", controllers.GetSalesPerformance)
   ```

2. **Implement controller** in `controllers/sales_controller.go`:
   ```go
   func GetSalesPerformance(c echo.Context) error {
       type response struct { ... }
       var result response
       database.DB.Raw("SELECT ...").Scan(&result)
       return c.JSON(http.StatusOK, result)
   }
   ```

3. **Model Updates**: add structs to `models/sales.go` if new tables/columns introduced.

4. **Tests**: place unit tests under `controllers/..._test.go` using `httptest`.

---

## 5. Database Access Patterns

- Use `database.DB.Raw()` for complex aggregations and analytics queries.
- For typical CRUD, leverage `database.DB.Model(&Model{}).Where().Find()`.
- Always sanitize inputs; parameterize raw SQL to avoid injection.
- When interacting with large datasets, paginate results with `Limit/Offset` or stream.

### Migrations

Schema changes live in `admin-cms/database/migrations/`. Coordinate with Laravel team to keep backend structs in sync.

---

## 6. Interacting with the AI Service

- Endpoints under `/api/ai/*` act as pass-through + persistence.
- Functions like `PredictSalesSummary`, `GetInventoryPrescription`, `GetFinancialPrescription`, `AnalyzeAnomaliesWithActions`, and `RunScenarioSimulation` demonstrate the pattern:
  1. Prepare payload (dates, filters).
  2. Call FastAPI via `postToAIService`.
  3. Persist recommendations in `models.PrescriptiveRecommendation` or `ScenarioSimulation` when meaningful.
  4. Return AI response to clients.

Ensure `AI_SERVICE_URL` is configured and reachable before deploying features relying on these paths.

---

## 7. Security & Middleware

- JWT validation occurs in `middleware/AuthMiddleware`. Update token claims in `utils/jwt.go` when RBAC changes.
- `middleware/security.go` sets headers (CSP, HSTS, rate limiting). Extend as needed.
- Always respect user role context when returning data (e.g., restrict executive-only dashboards).

---

## 8. Testing & Quality

| Task                 | Command |
|----------------------|---------|
| Run all tests        | `go test ./...` |
| Lint (golangci-lint) | `golangci-lint run` |
| Vet                  | `go vet ./...` |

For integration tests, stand up a local MySQL instance and use fixtures. Mock AI service by overriding `AI_SERVICE_URL` to hit httptest server.

---

## 9. Deployment Pipeline

1. Build artifact: `go build -o bin/backend-api ./...`.
2. Containerize using multi-stage Dockerfile (builder + scratch/alpine).
3. Deploy via CI (GitHub Actions) to Kubernetes / ECS with environment-specific configuration.
4. Run database migrations separately (Laravel artisan).

---

## 10. Troubleshooting Cheatsheet

| Symptom                           | Possible Fix |
|-----------------------------------|--------------|
| 401 Unauthorized                  | Check Authorization header formatting, token expiry |
| 500 AI prediction failed          | Inspect FastAPI logs, ensure SQL dataset available |
| DB deadlock errors                | Reduce transaction scope, add retries |
| Slow analytics queries            | Verify indexes (see migrations), cache results |

Keep this guide updated when introducing new modules, event-driven patterns, or breaking architectural changes.
