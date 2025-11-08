# AI Service (FastAPI)

Python microservice delivering predictive forecasts, anomaly detection, prescriptive recommendations, and what-if simulations for iDash.

---

## ✨ Capabilities

- Prophet-based forecasts for Sales, Production, Finance, Inventory.
- Grouped sales prediction summaries with action recommendations.
- Inventory optimization (reorder points, safety stock) and slow-mover detection.
- Financial prescriptive insights (expense controls, surplus deployment).
- Isolation Forest anomaly detection with severity scoring and recommended actions.
- Scenario simulations modeling price/volume/cost adjustments.

---

## ⚙️ Setup

```bash
cd ai-service
python -m venv venv
source venv/bin/activate  # Windows: venv\Scripts\activate
pip install -r requirements.txt

cp .env.example .env  # configure DB_* values
uvicorn main:app --reload --port 8000
```

Expose Swagger UI at `http://localhost:8000/docs` and ReDoc at `/redoc`.

---

## 📡 Key Endpoints

| Method | Path                              | Description |
|--------|-----------------------------------|-------------|
| POST   | `/api/forecast/sales`             | Sales forecast (Prophet) |
| POST   | `/api/forecast/production`        | Production forecast |
| POST   | `/api/forecast/finance`           | Financial forecast |
| POST   | `/api/forecast/inventory`         | Inventory forecast |
| POST   | `/api/predict/sales/summary`      | Next-period sales prediction + recommendations |
| POST   | `/api/prescribe/inventory`        | Inventory optimisation + slow movers |
| POST   | `/api/prescribe/finance`          | Financial prescriptive actions |
| POST   | `/api/analyze`                    | Trend, anomaly, correlation analysis |
| POST   | `/api/analyze/anomaly/enriched`   | Enriched anomaly feed with recommended actions |
| POST   | `/api/scenario/whatif`            | Scenario impact calculation |
| GET    | `/api/forecast/{forecast_id}`     | Retrieve saved forecast metadata |
| GET    | `/health`                         | Health probe |

All payloads use JSON; see `main.py` Pydantic models for schemas.

---

## 🧭 Development Workflow

1. Extend database helper (`database/database.py`) when new datasets are required.
2. Add business logic in the appropriate service module (`forecast_service`, `analysis_service`, `prescriptive_service`).
3. Expose new routes in `main.py` with Pydantic request/response models.
4. Coordinate with the Go backend to proxy new functionality to frontend clients.

Detailed patterns and testing guidance are documented in the [AI Service Developer Guide](../docs/developer/AI_SERVICE_GUIDE.md).

---

## 🧪 Testing & Quality

| Task            | Command |
|-----------------|---------|
| Unit tests      | `pytest` |
| Formatting      | `black .` |
| Linting         | `flake8` or `ruff .` |

Mock database methods during unit tests to avoid hitting live MySQL.

---

## 🚀 Deployment Notes

- Containerize with slim Python base and install system deps required by Prophet.
- Run under Gunicorn + Uvicorn workers for production (`gunicorn -k uvicorn.workers.UvicornWorker main:app -w 2 --threads 4`).
- Configure observability (Prometheus metrics, structured logging) in alignment with the [SRE Runbook](../docs/SRE_RUNBOOK.md).

---

## 📎 Resources

- [Setup Guide](../docs/SETUP_GUIDE.md)
- [Feature Catalog](../docs/FEATURE_CATALOG.md)
- [Sample Data Guide](../docs/SAMPLE_DATA_GUIDE.md)

