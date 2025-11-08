# AI-Based Integrated Management Dashboard

A multi-service enterprise intelligence platform delivering predictive and prescriptive analytics across sales, production, finance, inventory, HR, and supply chain.

---

## 🧭 Overview

- **Frontend** (React + Vite + Material UI) renders dynamic dashboards backed by AI insights.
- **Backend API** (Go + Echo + GORM) exposes secured REST endpoints and orchestrates AI calls.
- **Admin CMS** (Laravel) manages data ingestion, configuration, and demo orchestration.
- **AI Service** (FastAPI + Prophet + scikit-learn) produces forecasts, anomaly detection, and recommendations.

The system is production-hardened with RBAC, caching, session-aware data fetching, and an extensible data model sized for 1000 crore turnover scenarios.

---

## 📁 Project Structure

```
AI based Dashboard/
├── admin-cms/        # Laravel admin portal & ETL
├── backend-api/      # Golang REST API
├── frontend/         # React dashboards
├── ai-service/       # FastAPI analytics microservice
├── docs/             # Documentation hub (setup, SRE, developer guides)
└── README.md
```

---

## 🚀 Feature Highlights

- **Role-based dashboards** with KPI widgets, trend charts, AI forecasts, anomaly alerts, and scenario impact cards.
- **Predictive intelligence** leveraging Prophet forecasts for sales, production, finance, and inventory.
- **Prescriptive recommendations** for sales recovery, inventory optimization, and spend control.
- **What-if simulations** to model price, volume, and cost changes with incremental profit projections.
- **Automated demo data** generator covering all modules while preserving user accounts.
- **Enterprise security** via JWT auth, security headers, rate limiting, and audit-ready logs.

Refer to the [Feature Catalog](./docs/FEATURE_CATALOG.md) for a detailed module-by-module breakdown.

---

## 🛠️ Quick Start

1. Provision prerequisites (Node 18, Go 1.21, PHP 8.1, Python 3.10, MySQL 8.0).
2. Follow the [Setup Guide](./docs/SETUP_GUIDE.md) for environment variables, migrations, and service bootstrapping.
3. Generate demo data via `php artisan idash:demo-data --years=2 --refresh` (see [Sample Data Guide](./docs/SAMPLE_DATA_GUIDE.md)).
4. Start services:
   ```bash
   # Frontend
   cd frontend && npm install && npm run dev

   # Backend API
   cd backend-api && go run main.go

   # Admin CMS
   cd admin-cms && composer install && php artisan serve

   # AI Service
   cd ai-service && pip install -r requirements.txt && uvicorn main:app --reload
   ```

---

## 🧑‍💻 Developer Resources

| Area        | Guide |
|-------------|-------|
| Frontend    | [React Developer Guide](./docs/developer/FRONTEND_GUIDE.md) |
| Backend API | [Go Developer Guide](./docs/developer/BACKEND_API_GUIDE.md) |
| Admin CMS   | [Laravel Developer Guide](./docs/developer/ADMIN_CMS_GUIDE.md) |
| AI Service  | [FastAPI Developer Guide](./docs/developer/AI_SERVICE_GUIDE.md) |
| Operations  | [SRE Runbook](./docs/SRE_RUNBOOK.md) |

Additional references:
- [Marketing Brief](./docs/MARKETING_BRIEF.md)
- [Feature Catalog](./docs/FEATURE_CATALOG.md)
- [Sample Data Guide](./docs/SAMPLE_DATA_GUIDE.md)

---

## 🔄 Deployment Notes

- Backend builds as a single static binary (`go build -o bin/backend-api`).
- Frontend bundles via `npm run build` and can be hosted on any static provider.
- Admin CMS runs on PHP-FPM with queue workers for ETL jobs.
- AI service ships as a containerized Uvicorn/Gunicorn stack.
- Consult the [SRE Runbook](./docs/SRE_RUNBOOK.md) for monitoring, scaling, and incident response.

---

## 🤝 Contributing

1. Fork the repository & create a feature branch.
2. Follow coding standards per service guide.
3. Ensure tests pass (`npm test`, `go test ./...`, `php artisan test`, `pytest`).
4. Submit a pull request referencing the associated issue / user story.

---

## 📄 License

MIT License © Development Team

### Windows Helper Scripts

- `setup.bat` – installs dependencies, copies `.env` templates, runs Laravel migrations, and prepares the AI service virtual environment.
- `run.bat` – launches backend API, admin CMS web server and queue worker, AI service, and frontend dev server, then opens the app in the browser.

